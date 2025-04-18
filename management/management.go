package management

import (
	"context"
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

// Applications returns an ApplicationService to manage applications.
func (c *Client) Applications() *ApplicationService {
	return &ApplicationService{client: c}
}
