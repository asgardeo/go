package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// APIResourceService handles API resource management operations.
type APIResourceService struct {
	client *Client
}

// APIResource represents an Asgardeo API resource matching APIResourceResponse.
type APIResource struct {
	ID                     string                          `json:"id,omitempty"`
	Name                   string                          `json:"name"`
	Identifier             string                          `json:"identifier"`
	Type                   string                          `json:"type,omitempty"`
	Description            string                          `json:"description,omitempty"`
	RequiresAuthorization  bool                            `json:"requiresAuthorization"`
	Scopes                 []ScopeGetModel                 `json:"scopes,omitempty"`
	SubscribedApplications []SubscribedApplicationGetModel `json:"subscribedApplications,omitempty"`
	Properties             []Property                      `json:"properties,omitempty"`
	Self                   string                          `json:"self,omitempty"`
}

// APIResourceCreateInput represents payload to create an API resource (matches APIResourceCreationModel schema).
type APIResourceCreateInput struct {
	Identifier            string               `json:"identifier"`
	Name                  string               `json:"name"`
	Description           string               `json:"description,omitempty"`
	RequiresAuthorization bool                 `json:"requiresAuthorization"`
	Scopes                []ScopeCreationModel `json:"scopes,omitempty"`
}

// APIResourceUpdateInput represents payload to update an API resource.
type APIResourceUpdateInput struct {
	Name        *string  `json:"name,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	Description *string  `json:"description,omitempty"`
}

// ScopeCreationModel represents a scope to be created or added to API resources.
type ScopeCreationModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// ScopeGetModel represents a retrieved scope with ID and metadata.
type ScopeGetModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// APIResourcePatchModel represents the fields for patching an API resource.
type APIResourcePatchModel struct {
	Name          *string              `json:"name,omitempty"`
	Description   *string              `json:"description,omitempty"`
	AddedScopes   []ScopeCreationModel `json:"addedScopes,omitempty"`
	RemovedScopes []string             `json:"removedScopes,omitempty"`
}

// SubscribedApplicationGetModel represents an application subscribed to an API resource.
type SubscribedApplicationGetModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Property represents a custom property of an API resource.
type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// APIResources returns an APIResourceService.
func (c *Client) APIResources() *APIResourceService {
	return &APIResourceService{client: c}
}

// List retrieves all API resources.
func (s *APIResourceService) List(ctx context.Context) ([]APIResource, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources", s.client.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Resources []APIResource `json:"resources"`
	}
	if err := s.client.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return resp.Resources, nil
}

// Get retrieves an API resource by its ID.
func (s *APIResourceService) Get(ctx context.Context, id string) (*APIResource, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var r APIResource
	if err := s.client.doRequest(req, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Create creates a new API resource.
func (s *APIResourceService) Create(ctx context.Context, input APIResourceCreateInput) (*APIResource, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources", s.client.baseURL)
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var created APIResource
	if err := s.client.doRequest(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates an existing API resource.
func (s *APIResourceService) Update(ctx context.Context, id string, input APIResourceUpdateInput) (*APIResource, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s", s.client.baseURL, id)
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var updated APIResource
	if err := s.client.doRequest(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Patch applies partial updates to an API resource using JSON Patch semantics (limited to name, description, scopes).
func (s *APIResourceService) Patch(ctx context.Context, id string, input APIResourcePatchModel) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s", s.client.baseURL, id)
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// Delete removes an API resource by its ID.
func (s *APIResourceService) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// GetScopes lists the scopes for a specific API resource.
func (s *APIResourceService) GetScopes(ctx context.Context, id string) ([]ScopeGetModel, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s/scopes", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var scopes []ScopeGetModel
	if err := s.client.doRequest(req, &scopes); err != nil {
		return nil, err
	}
	return scopes, nil
}

// AddScopes replaces scopes for a specific API resource (PUT) based on spec.
func (s *APIResourceService) AddScopes(ctx context.Context, id string, scopes []ScopeCreationModel) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s/scopes", s.client.baseURL, id)
	body, err := json.Marshal(scopes)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	// spec returns 204 No Content
	return s.client.doRequest(req, nil)
}

// DeleteScope deletes a single scope from an API resource.
func (s *APIResourceService) DeleteScope(ctx context.Context, id, scopeName string) error {
	endpoint := fmt.Sprintf("%s/api/server/v1/api-resources/%s/scopes/%s", s.client.baseURL, id, scopeName)
	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// ListScopes lists all scopes in the tenant (global scopes).
func (s *APIResourceService) ListScopes(ctx context.Context, filter string) ([]ScopeGetModel, error) {
	// optional filter param scopeFilter
	u := fmt.Sprintf("%s/scopes", s.client.baseURL)
	if filter != "" {
		u = fmt.Sprintf("%s?scopeFilter=%s", u, url.QueryEscape(filter))
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}
	var scopes []ScopeGetModel
	if err := s.client.doRequest(req, &scopes); err != nil {
		return nil, err
	}
	return scopes, nil
}
