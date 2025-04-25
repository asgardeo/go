package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Organization represents an Asgardeo organization.
type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`
	// Add other fields as needed
}

// OrganizationCreateInput represents payload to create an organization.
type OrganizationCreateInput struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
	Type        string `json:"type,omitempty"`
	// Add other fields as needed
}

// ListOrganizationsParams holds optional query parameters for listing organizations.
type ListOrganizationsParams struct {
	Limit     int
	Filter    string
	Recursive bool
}

// toQuery builds URL query string from params.
func (p *ListOrganizationsParams) toQuery() string {
	vals := url.Values{}
	if p.Limit > 0 {
		vals.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Filter != "" {
		vals.Set("filter", p.Filter)
	}
	vals.Set("recursive", fmt.Sprintf("%v", p.Recursive))
	qs := vals.Encode()
	if qs != "" {
		return "?" + qs
	}
	return ""
}

// OrganizationService handles organization management operations.
type OrganizationService struct {
	client *Client
}

// ListOrganizationsResponse models the response for listing organizations.
type ListOrganizationsResponse struct {
	Count         int            `json:"count"`
	TotalResults  int            `json:"totalResults"`
	Organizations []Organization `json:"organizations"`
}

// List retrieves a list of organizations with optional parameters.
func (s *OrganizationService) List(ctx context.Context, params *ListOrganizationsParams) (*ListOrganizationsResponse, error) {
	path := "/api/server/v1/organizations"
	if params != nil {
		path += params.toQuery()
	}
	endpoint := fmt.Sprintf("%s%s", s.client.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var resp ListOrganizationsResponse
	if err := s.client.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves an organization by its ID.
func (s *OrganizationService) Get(ctx context.Context, id string) (*Organization, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/organizations/%s", s.client.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	var org Organization
	if err := s.client.doRequest(req, &org); err != nil {
		return nil, err
	}
	return &org, nil
}

// Create creates a new organization.
func (s *OrganizationService) Create(ctx context.Context, input OrganizationCreateInput) (*Organization, error) {
	endpoint := fmt.Sprintf("%s/api/server/v1/organizations", s.client.baseURL)
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var created Organization
	if err := s.client.doRequest(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}
