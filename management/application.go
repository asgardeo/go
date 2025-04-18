package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ApplicationService handles application management operations.
type ApplicationService struct {
	client *Client
}

// Application represents an Asgardeo application.
type Application struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Add other fields as needed
}

// AdvancedConfigurations defines skip consent behaviors.
type AdvancedConfigurations struct {
	SkipLogoutConsent bool `json:"skipLogoutConsent"`
	SkipLoginConsent  bool `json:"skipLoginConsent"`
}

// AssociatedRoles defines roles association for application.
type AssociatedRoles struct {
	AllowedAudience string   `json:"allowedAudience"`
	Roles           []string `json:"roles"`
}

// InboundOIDCConfig holds minimal OIDC inbound protocol settings.
type InboundOIDCConfig struct {
	GrantTypes        []string `json:"grantTypes"`
	IsFAPIApplication bool     `json:"isFAPIApplication"`
}

// InboundProtocolConfiguration wraps protocol configs.
type InboundProtocolConfiguration struct {
	OIDC *InboundOIDCConfig `json:"oidc"`
}

// ApplicationCreateInput represents payload to create an application.
type ApplicationCreateInput struct {
	Name                         string                        `json:"name"`
	AdvancedConfigurations       *AdvancedConfigurations       `json:"advancedConfigurations"`
	TemplateID                   string                        `json:"templateId"`
	AssociatedRoles              *AssociatedRoles              `json:"associatedRoles"`
	InboundProtocolConfiguration *InboundProtocolConfiguration `json:"inboundProtocolConfiguration"`
}

// ListApplicationsResponse models the response for listing applications.
type ListApplicationsResponse struct {
	Count        int           `json:"count"`
	TotalResults int           `json:"totalResults"`
	Applications []Application `json:"applications"`
}

// ListApplicationsParams holds optional query parameters for listing applications.
type ListApplicationsParams struct {
	Limit  int    // max number of results
	Offset int    // starting index
	Filter string // filter expression
	Sort   string // sort order
}

// toQuery builds URL query string from params.
func (p *ListApplicationsParams) toQuery() string {
	vals := url.Values{}
	if p.Limit > 0 {
		vals.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Offset > 0 {
		vals.Set("offset", fmt.Sprintf("%d", p.Offset))
	}
	if p.Filter != "" {
		vals.Set("filter", p.Filter)
	}
	if p.Sort != "" {
		vals.Set("sort", p.Sort)
	}
	qs := vals.Encode()
	if qs != "" {
		return "?" + qs
	}
	return ""
}

// List retrieves a list of applications with optional parameters.
func (s *ApplicationService) List(ctx context.Context, params *ListApplicationsParams) (*ListApplicationsResponse, error) {
	path := "/api/server/v1/applications"
	if params != nil {
		path += params.toQuery()
	}
	endpoint := fmt.Sprintf("%s%s", s.client.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var resp ListApplicationsResponse
	if err := s.client.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an application by its ID.
func (s *ApplicationService) Get(ctx context.Context, id string) (*Application, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var app Application
	if err := s.client.doRequest(req, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

// Create creates a new application.
func (s *ApplicationService) Create(ctx context.Context, input ApplicationCreateInput) (*Application, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications", s.client.baseURL)
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var created Application
	if err := s.client.doRequest(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// PatchOperation represents a JSON Patch (RFC6902) operation.
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// Update updates an existing application (full replace via PUT).
func (s *ApplicationService) Update(ctx context.Context, id string, app Application) (*Application, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s", s.client.baseURL, id)
	body, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var updated Application
	if err := s.client.doRequest(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes an application by its ID.
func (s *ApplicationService) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// RegenerateClientSecret rotates the client secret for an application.
func (s *ApplicationService) RegenerateClientSecret(ctx context.Context, id string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/client-secret", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return "", err
	}
	var resp struct {
		Secret string `json:"secret"`
	}
	if err := s.client.doRequest(req, &resp); err != nil {
		return "", err
	}
	return resp.Secret, nil
}

// GetCertificate retrieves the PEM certificate for an application.
func (s *ApplicationService) GetCertificate(ctx context.Context, id string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/certificate", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	var resp struct {
		Certificate string `json:"certificate"`
	}
	if err := s.client.doRequest(req, &resp); err != nil {
		return "", err
	}
	return resp.Certificate, nil
}

// UpdateCertificate updates the PEM certificate for an application.
func (s *ApplicationService) UpdateCertificate(ctx context.Context, id, certificate string) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/certificate", s.client.baseURL, id)
	body, err := json.Marshal(map[string]string{"certificate": certificate})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// OIDCConfig represents the OIDC inbound protocol configuration.
type OIDCConfig struct {
	ClientID                string   `json:"clientId"`
	ClientSecret            string   `json:"clientSecret"`
	GrantTypes              []string `json:"grantTypes"`
	ResponseTypes           []string `json:"responseTypes"`
	CallbackURLs            []string `json:"callbackUrls"`
	LogoutURLs              []string `json:"logoutUrls"`
	TokenEndpointAuthMethod string   `json:"tokenEndpointAuthMethod"`
}

// GetOIDCConfig retrieves the OIDC configuration for an application.
func (s *ApplicationService) GetOIDCConfig(ctx context.Context, id string) (*OIDCConfig, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/inbound-protocols/oidc", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var cfg OIDCConfig
	if err := s.client.doRequest(req, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// UpdateOIDCConfig updates the OIDC configuration for an application.
func (s *ApplicationService) UpdateOIDCConfig(ctx context.Context, id string, cfg OIDCConfig) (*OIDCConfig, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/inbound-protocols/oidc", s.client.baseURL, id)
	body, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var updated OIDCConfig
	if err := s.client.doRequest(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// SAMLConfig represents the SAML inbound protocol configuration.
type SAMLConfig struct {
	Issuer                string   `json:"issuer"`
	AssertionConsumerUrls []string `json:"assertionConsumerUrls"`
	SingleLogoutUrls      []string `json:"singleLogoutUrls"`
	SigningEnabled        bool     `json:"signingEnabled"`
	EncryptAssertion      bool     `json:"encryptAssertion"`
	EnableSso             bool     `json:"enableSso"`
}

// GetSAMLConfig retrieves the SAML configuration for an application.
func (s *ApplicationService) GetSAMLConfig(ctx context.Context, id string) (*SAMLConfig, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/inbound-protocols/saml2", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var cfg SAMLConfig
	if err := s.client.doRequest(req, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// UpdateSAMLConfig updates the SAML configuration for an application.
func (s *ApplicationService) UpdateSAMLConfig(ctx context.Context, id string, cfg SAMLConfig) (*SAMLConfig, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/inbound-protocols/saml2", s.client.baseURL, id)
	body, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var updated SAMLConfig
	if err := s.client.doRequest(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// ApplicationTemplateService handles application template operations.
type ApplicationTemplateService struct {
	client *Client
}

// ApplicationTemplate represents an application template.
type ApplicationTemplate struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Templates returns an ApplicationTemplateService.
func (c *Client) Templates() *ApplicationTemplateService {
	return &ApplicationTemplateService{client: c}
}

// List retrieves all application templates.
func (s *ApplicationTemplateService) List(ctx context.Context) ([]ApplicationTemplate, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/templates", s.client.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Templates []ApplicationTemplate `json:"templates"`
	}
	if err := s.client.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return resp.Templates, nil
}

// Get retrieves a single application template by ID.
func (s *ApplicationTemplateService) Get(ctx context.Context, id string) (*ApplicationTemplate, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/templates/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var tpl ApplicationTemplate
	if err := s.client.doRequest(req, &tpl); err != nil {
		return nil, err
	}
	return &tpl, nil
}

// SharedApplicationService handles shared application operations.
type SharedApplicationService struct {
	client *Client
}

// SharedApplications returns a SharedApplicationService.
func (c *Client) SharedApplications() *SharedApplicationService {
	return &SharedApplicationService{client: c}
}

// List retrieves shared applications.
func (s *SharedApplicationService) List(ctx context.Context, params *ListApplicationsParams) (*ListApplicationsResponse, error) {
	path := "/api/server/v1/applications/shared"
	if params != nil {
		path += params.toQuery()
	}
	endpoint := fmt.Sprintf("%s%s", s.client.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var resp ListApplicationsResponse
	if err := s.client.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a shared application by its ID.
func (s *SharedApplicationService) Get(ctx context.Context, id string) (*Application, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/shared/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var app Application
	if err := s.client.doRequest(req, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

// AuthorizedAPI defines an API authorized for the application.
type AuthorizedAPI struct {
	APIID  string   `json:"apiId"`
	Scopes []string `json:"scopes,omitempty"`
}

// AuthorizedAPIsService handles operations on APIs authorized for an application.
type AuthorizedAPIsService struct {
	client *Client
	appID  string
}

// AuthorizedAPIs returns a AuthorizedAPIsService for the specified application.
func (c *Client) AuthorizedAPIs(appID string) *AuthorizedAPIsService {
	return &AuthorizedAPIsService{client: c, appID: appID}
}

// List retrieves the list of APIs authorized for the application.
func (s *AuthorizedAPIsService) List(ctx context.Context) ([]AuthorizedAPI, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/authorized-apis", s.client.baseURL, s.appID)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Apis []AuthorizedAPI `json:"apis"`
	}
	if err := s.client.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return resp.Apis, nil
}

// Get retrieves a specific authorized API by ID.
func (s *AuthorizedAPIsService) Get(ctx context.Context, apiID string) (*AuthorizedAPI, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/authorized-apis/%s", s.client.baseURL, s.appID, apiID)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var api AuthorizedAPI
	if err := s.client.doRequest(req, &api); err != nil {
		return nil, err
	}
	return &api, nil
}

// Update replaces the list of authorized APIs for the application.
func (s *AuthorizedAPIsService) Update(ctx context.Context, apis []AuthorizedAPI) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/authorized-apis", s.client.baseURL, s.appID)
	body, err := json.Marshal(map[string][]AuthorizedAPI{"apis": apis})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// Delete removes an authorized API from the application.
func (s *AuthorizedAPIsService) Delete(ctx context.Context, apiID string) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/applications/%s/authorized-apis/%s", s.client.baseURL, s.appID, apiID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}
