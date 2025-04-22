package management

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is the Asgardeo management API client.
type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	token        string
	ctx          context.Context // context for requests
}

// Option configures the management client.
type Option func(*Client) error

// New initializes a new Asgardeo management client with baseURL and options.
func New(baseURL string, opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: http.DefaultClient,
		ctx:        context.Background(),
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithClientCredentials sets up client credentials (client_credentials grant) and context.
func WithClientCredentials(ctx context.Context, clientID, clientSecret string) Option {
	return func(c *Client) error {
		c.ctx = ctx
		c.clientID = clientID
		c.clientSecret = clientSecret
		return nil
	}
}

// WithStaticToken configures a static bearer token.
func WithStaticToken(token string) Option {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

// WithHTTPClient configures a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		c.httpClient = hc
		return nil
	}
}

// authenticate performs OAuth2 client credentials grant to retrieve an access token.
func (c *Client) authenticate() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("scope", "SYSTEM")
	endpoint := c.baseURL + "/oauth2/token"
	resp, err := c.httpClient.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: %s", body)
	}
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.AccessToken == "" {
		return errors.New("received empty access token")
	}
	c.token = result.AccessToken
	return nil
}

// exchangeTokenForSubOrg exchanges the parent organization token for a sub-organization token
func (c *Client) exchangeTokenForSubOrg(orgID string) (string, error) {
	if c.token == "" {
		if err := c.authenticate(); err != nil {
			return "", fmt.Errorf("failed to authenticate to get base token for exchange: %w", err)
		}
	}

	if orgID == "" {
		return "", errors.New("organization ID cannot be empty")
	}

	data := url.Values{}
	data.Set("grant_type", "organization_switch")
	data.Set("scope", "SYSTEM")
	data.Set("switching_organization", orgID)
	data.Set("token", c.token)

	endpoint := c.baseURL + "/oauth2/token"
	req, err := http.NewRequestWithContext(c.ctx, "POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token exchange request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Add Basic Authorization header with base64(clientId:clientSecret)
	creds := c.clientID + ":" + c.clientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(creds))
	req.Header.Set("Authorization", "Basic "+basicAuth)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, body)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token exchange response: %w", err)
	}

	if result.AccessToken == "" {
		return "", errors.New("received empty access token from token exchange")
	}

	return result.AccessToken, nil
}

// doRequest wraps HTTP requests with authentication and JSON handling.
func (c *Client) doRequest(req *http.Request, v interface{}) error {
	if c.token == "" {
		if err := c.authenticate(); err != nil {
			return err
		}
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	// Skip decoding empty body
	if v != nil {
		// if no body, skip
		if resp.ContentLength == 0 {
			return nil
		}
		// attempt decode; ignore EOF
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		return nil
	}
	return nil
}

// doSubOrgRequest handles requests to sub-organization resources with token exchange
func (c *Client) doSubOrgRequest(orgID string, method, path string, body interface{}, v interface{}) error {
	// Exchange parent org token for sub-org token
	subOrgToken, err := c.exchangeTokenForSubOrg(orgID)
	if err != nil {
		return err
	}

	// Ensure path starts with a slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Sub-org requests include /o/ in the path
	fullURL := c.baseURL + "/o" + path

	// Prepare the request body if provided
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create the request
	req, err := http.NewRequestWithContext(c.ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers with sub-org token
	req.Header.Set("Authorization", "Bearer "+subOrgToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http client error: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	// Skip decoding empty body
	if v != nil {
		// if no body, skip
		if resp.ContentLength == 0 {
			return nil
		}
		// attempt decode; ignore EOF
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

// Applications returns an ApplicationService to manage applications.
func (c *Client) Applications() *ApplicationService {
	return &ApplicationService{client: c}
}

// SubOrgApplications returns a SubOrgApplicationService to manage applications in a sub-organization
func (c *Client) SubOrgApplications() *SubOrgApplicationService {
	return &SubOrgApplicationService{client: c}
}

func (c *Client) Organizations() *OrganizationService {
	return &OrganizationService{client: c}
}
