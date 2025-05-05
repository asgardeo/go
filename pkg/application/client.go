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

func (c *ApplicationClient) AuthorizeAPI(ctx context.Context, appID string, authorizedAPI AddAuthorizedAPIJSONRequestBody) (*http.Response, error) {
	resp, err := c.apiClient.AddAuthorizedAPIWithResponse(ctx, appID, authorizedAPI)
	if err != nil {
		return nil, fmt.Errorf("Failed to authorize api: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to authorize api: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return &http.Response{}, nil
}
