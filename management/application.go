package management

import (
	"context"
	"fmt"
	"net/http"

	applications "github.com/asgardeo/go/management/applications"
)

// ApplicationService handles application management operations.
type ApplicationService struct {
	client applications.ClientWithResponsesInterface
}

type ApplicationCreateInput struct {
	Name string `json:"name"`
}

// List retrieves a list of applications with optional parameters.
func (s *ApplicationService) List(ctx context.Context, params *applications.GetAllApplicationsParams) (*applications.ApplicationListResponse, error) {

	resp, err := s.client.GetAllApplicationsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to list applications: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// Get retrieves an application by its ID.
func (s *ApplicationService) Get(ctx context.Context, id string) (*applications.ApplicationResponseModel, error) {

	resp, err := s.client.GetApplicationWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get application: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// Create creates a new application.
func (s *ApplicationService) Create(ctx context.Context, application applications.ApplicationModel) error {

	resp, err := s.client.CreateApplicationWithResponse(ctx, nil, application)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to create application: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return nil
}

// Delete removes an application by its ID.
func (s *ApplicationService) Delete(ctx context.Context, id string) error {

	resp, err := s.client.DeleteApplicationWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed to delete application: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return nil
}
