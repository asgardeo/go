/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// AuthMethod represents the authentication method to use
type AuthMethod string

const (
	// AuthMethodToken represents direct token authentication
	AuthMethodToken AuthMethod = "token"
	// AuthMethodClientCredentials represents client credentials grant type
	AuthMethodClientCredentials AuthMethod = "client_credentials"
)

// TokenResponse represents the response from the OAuth token endpoint
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// ClientConfig contains the configuration for the API clients
type ClientConfig struct {
	// BaseURL is the base URL for the API
	BaseURL string
	// HTTPClient is the HTTP client to use for requests
	HTTPClient *http.Client
	// APIKey is the API key to use for authentication (deprecated, use Token instead)
	APIKey string
	// Timeout is the timeout for requests
	Timeout time.Duration

	// Auth related fields
	// AuthMethod is the authentication method to use
	AuthMethod AuthMethod
	// Token is the direct authentication token to use if AuthMethod is AuthMethodToken
	Token string
	// OAuth2ClientID is the client ID to use for client credentials grant type
	OAuth2ClientID string
	// OAuth2ClientSecret is the client secret to use for client credentials grant type
	OAuth2ClientSecret string

	// Token cache
	currentToken      string
	tokenExpiresAt    time.Time
	tokenRefreshMutex sync.Mutex
}

// DefaultClientConfig returns a default configuration for the API clients
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
}

// WithBaseURL sets the base URL for the API
func (c *ClientConfig) WithBaseURL(baseURL string) *ClientConfig {
	c.BaseURL = baseURL
	return c
}

// WithHTTPClient sets the HTTP client to use for requests
func (c *ClientConfig) WithHTTPClient(httpClient *http.Client) *ClientConfig {
	c.HTTPClient = httpClient
	return c
}

// WithAPIKey sets the API key to use for authentication (deprecated, use WithToken instead)
func (c *ClientConfig) WithAPIKey(apiKey string) *ClientConfig {
	c.APIKey = apiKey
	return c
}

// WithTimeout sets the timeout for requests
func (c *ClientConfig) WithTimeout(timeout time.Duration) *ClientConfig {
	c.Timeout = timeout
	if c.HTTPClient != nil {
		c.HTTPClient.Timeout = timeout
	}
	return c
}

// WithToken sets a static token for authentication
func (c *ClientConfig) WithToken(token string) *ClientConfig {
	c.AuthMethod = AuthMethodToken
	c.Token = token
	return c
}

// WithClientCredentials sets up client credentials grant type authentication
func (c *ClientConfig) WithClientCredentials(clientID, clientSecret string) *ClientConfig {
	c.AuthMethod = AuthMethodClientCredentials
	c.OAuth2ClientID = clientID
	c.OAuth2ClientSecret = clientSecret
	return c
}

// GetToken returns the current valid token, fetching a new one if necessary
func (c *ClientConfig) GetToken(ctx context.Context) (string, error) {
	c.tokenRefreshMutex.Lock()
	defer c.tokenRefreshMutex.Unlock()

	// If using direct token method and we have a token, return it
	if c.AuthMethod == AuthMethodToken && c.Token != "" {
		return c.Token, nil
	}

	// If using client credentials and we have a cached token that's not expired, return it
	if c.AuthMethod == AuthMethodClientCredentials {
		// If we have a valid token that's not expired (with 30s buffer), return it
		if c.currentToken != "" && time.Now().Add(30*time.Second).Before(c.tokenExpiresAt) {
			return c.currentToken, nil
		}

		// Otherwise, fetch a new token
		token, expiresIn, err := c.fetchClientCredentialsToken(ctx)
		if err != nil {
			return "", err
		}

		// Cache the token
		c.currentToken = token
		c.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
		return token, nil
	}

	// If we get here, we don't have a valid authentication method
	return "", fmt.Errorf("no valid authentication method configured")
}

// fetchClientCredentialsToken fetches a new token using client credentials grant type
func (c *ClientConfig) fetchClientCredentialsToken(ctx context.Context) (string, int, error) {
	if c.OAuth2ClientID == "" || c.OAuth2ClientSecret == "" {
		return "", 0, fmt.Errorf("client ID, client secret, and token URL are required for client credentials grant type")
	}

	// Prepare request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "SYSTEM") // Request system scope

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create token request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.OAuth2ClientID, c.OAuth2ClientSecret)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("failed to obtain token: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", 0, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", 0, fmt.Errorf("received empty access token")
	}

	return tokenResp.AccessToken, tokenResp.ExpiresIn, nil
}
