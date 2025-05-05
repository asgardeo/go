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

func (c *APIResourceClient) List(ctx context.Context, params *GetAPIResourcesParams) (*APIResourceListResponse, error) {
	resp, err := c.apiClient.GetAPIResourcesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Failed to list api resources: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list api resources: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (c *APIResourceClient) Get(ctx context.Context, id string) (*APIResourceResponse, error) {
	resp, err := c.apiClient.GetApiResourcesApiResourceIdWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get api resource: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to get api resource: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (c *APIResourceClient) GetByName(ctx context.Context, name string) (*[]APIResourceListItem, error) {
	filter := Filter(fmt.Sprintf("name eq %s", name))
	params := GetAPIResourcesParams{
		Filter: &filter,
	}
	resp, err := c.apiClient.GetAPIResourcesWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("Failed to list api resources: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list api resources: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	// Get the list of API resources from the response.
	apiResources := resp.JSON200.APIResources
	return apiResources, nil
}

func (c *APIResourceClient) GetByIdentifier(ctx context.Context, identifier string) (*APIResourceListItem, error) {
	filter := Filter(fmt.Sprintf("identifier eq %s", identifier))
	params := GetAPIResourcesParams{
		Filter: &filter,
	}
	resp, err := c.apiClient.GetAPIResourcesWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("Failed to get api resource: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to get api resource: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	// Since the identifier is unique, we can return the first item in the list.
	if len(*resp.JSON200.APIResources) == 0 {
		return nil, fmt.Errorf("No API resource found with identifier: %s", identifier)
	}
	return &(*resp.JSON200.APIResources)[0], nil
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
