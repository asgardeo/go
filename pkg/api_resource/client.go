package api_resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

// APIResourceClient is a wrapper around the generated client for the API Resource Management API.
type APIResourceClient struct {
	config    *config.ClientConfig
	apiClient *ClientWithResponses
}

// Creates a new API Resource Management API client.
func New(cfg *config.ClientConfig) (*APIResourceClient, error) {

	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiClient, err := NewClientWithResponses(
		cfg.BaseURL+"/api/server/v1",
		WithHTTPClient(cfg.HTTPClient),
		WithRequestEditorFn(typedAuthEditorFn),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create api resource client: %w", err)
	}

	return &APIResourceClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *APIResourceClient) List(ctx context.Context, limit *int, before *string, after *string, filter *string) (*APIResourceListResponse, error) {
	params := GetAPIResourcesParams{
		Limit:  limit,
		Before: before,
		After:  after,
		Filter: filter,
	}
	resp, err := c.apiClient.GetAPIResourcesWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("Failed to list api resources: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list api resources: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (c *APIResourceClient) Create(ctx context.Context, apiResource *AddAPIResourceJSONRequestBody) (*AddAPIResourceResponse, error) {
	resp, err := c.apiClient.AddAPIResourceWithResponse(ctx, *apiResource)
	if err != nil {
		return nil, fmt.Errorf("Failed to create api resource: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create api resource: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp, nil
}
