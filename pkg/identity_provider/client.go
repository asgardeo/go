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

package identity_provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

type IdentityProviderClient struct {
	config    *config.ClientConfig
	apiClient *ClientWithResponses
}

func New(cfg *config.ClientConfig) (*IdentityProviderClient, error) {
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
		return nil, fmt.Errorf("Failed to create identity provider client: %w", err)
	}

	return &IdentityProviderClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *IdentityProviderClient) List(ctx context.Context, idpGetParams *GetIDPsParams) (*IdentityProviderListResponse, error) {
	resp, err := c.apiClient.GetIDPsWithResponse(ctx, idpGetParams)
	if err != nil {
		return nil, fmt.Errorf("Failed to list identity providers: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list identity providers: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}

	return resp.JSON200, nil
}
