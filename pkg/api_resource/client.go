/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package api_resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/api_resource/internal"
	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

// APIResourceClient is a wrapper around the generated client for the API Resource Management API.
type APIResourceClient struct {
	config            *config.ClientConfig
	apiResourceClient *internal.ClientWithResponses
}

// Creates a new API Resource Management API client.
func New(cfg *config.ClientConfig) (*APIResourceClient, error) {

	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiResourceClient, err := internal.NewClientWithResponses(
		cfg.BaseURL+"/api/server/v1",
		internal.WithHTTPClient(cfg.HTTPClient),
		internal.WithRequestEditorFn(typedAuthEditorFn),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create api resource client: %w", err)
	}

	return &APIResourceClient{
		config:            cfg,
		apiResourceClient: apiResourceClient,
	}, nil
}

func (c *APIResourceClient) List(ctx context.Context, params *APIResourceListParamsModel) (*APIResourceListResponseModel, error) {
	resp, err := c.apiResourceClient.GetAPIResourcesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Failed to list api resources: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list api resources: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (c *APIResourceClient) Get(ctx context.Context, id string) (*APIResourceInfoResponseModel, error) {
	resp, err := c.apiResourceClient.GetApiResourcesApiResourceIdWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get api resource: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to get api resource: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (c *APIResourceClient) GetByName(ctx context.Context, name string) (*[]APIResourceListItemModel, error) {
	filter := internal.Filter(fmt.Sprintf("name eq %s", name))
	params := internal.GetAPIResourcesParams{
		Filter: &filter,
	}
	resp, err := c.apiResourceClient.GetAPIResourcesWithResponse(ctx, &params)
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

func (c *APIResourceClient) GetByIdentifier(ctx context.Context, identifier string) (*APIResourceListItemModel, error) {
	filter := internal.Filter(fmt.Sprintf("identifier eq %s", identifier))
	params := internal.GetAPIResourcesParams{
		Filter: &filter,
	}
	resp, err := c.apiResourceClient.GetAPIResourcesWithResponse(ctx, &params)
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

func (c *APIResourceClient) Create(ctx context.Context, apiResource *APIResourceCreateModel) (*APIResourceInfoResponseModel, error) {
	resp, err := c.apiResourceClient.AddAPIResourceWithResponse(ctx, *apiResource)
	if err != nil {
		return nil, fmt.Errorf("Failed to create api resource: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create api resource: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON201, nil
}
