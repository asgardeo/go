package management

import (
	"fmt"
	"net/http"
	"net/url"
)

// SubOrgApplicationService handles operations on applications within a sub-organization
type SubOrgApplicationService struct {
	client *Client
}

// GetAll retrieves all applications in the specified sub-organization
// It requires the organization ID of the sub-organization
// Parameters:
//   - orgID: The ID of the sub-organization
//   - queryParams: Optional query parameters like limit, offset, filter, etc.
func (s *SubOrgApplicationService) GetAll(orgID string, queryParams map[string]string) (*ListApplicationsResponse, error) {
	if orgID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	path := "/api/server/v1/applications"

	// Add query parameters if provided
	if len(queryParams) > 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, v)
		}
		path += "?" + params.Encode()
	}

	var appList ListApplicationsResponse
	err := s.client.doSubOrgRequest(orgID, http.MethodGet, path, nil, &appList)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications for organization %s: %w", orgID, err)
	}

	return &appList, nil
}

// GetByID retrieves a specific application in the specified sub-organization
// Parameters:
//   - orgID: The ID of the sub-organization
//   - appID: The ID of the application to retrieve
func (s *SubOrgApplicationService) GetByID(orgID string, appID string) (*Application, error) {
	if orgID == "" {
		return nil, fmt.Errorf("organization ID cannot be empty")
	}

	if appID == "" {
		return nil, fmt.Errorf("application ID cannot be empty")
	}

	path := fmt.Sprintf("/api/server/v1/applications/%s", appID)

	var app Application
	err := s.client.doSubOrgRequest(orgID, http.MethodGet, path, nil, &app)
	if err != nil {
		return nil, fmt.Errorf("failed to get application %s for organization %s: %w", appID, orgID, err)
	}

	return &app, nil
}
