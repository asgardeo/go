package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

const (
	defaultLimit  int = 10
	defaultOffset int = 0
)

// ApplicationClient is a wrapper around the generated client for the Application Management API
type ApplicationClient struct {
	config    *config.ClientConfig
	apiClient *ClientWithResponses
}

// NewWrapperClient creates a new Application Management API client
func New(cfg *config.ClientConfig) (*ApplicationClient, error) {

	// Create an auth request editor function
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	// Create a wrapper function that conforms to the RequestEditorFn type
	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		// Cast the generic function and call it
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	// Create the client with the auth editor
	apiClient, err := NewClientWithResponses(
		cfg.BaseURL+"/api/server/v1",
		WithHTTPClient(cfg.HTTPClient),
		WithRequestEditorFn(typedAuthEditorFn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create application client: %w", err)
	}

	return &ApplicationClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *ApplicationClient) List(ctx context.Context, limit, offset int) (*ApplicationListResponse, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = defaultOffset
	}

	params := GetAllApplicationsParams{
		Limit:  &limit,
		Offset: &offset,
	}
	resp, err := c.apiClient.GetAllApplicationsWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to list applications: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
