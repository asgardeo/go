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

package oidc_scope

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/oidc_scope/internal"
)

type OIDCScopeClient struct {
	config    *config.ClientConfig
	apiClient *internal.ClientWithResponses
}

func New(cfg *config.ClientConfig) (*OIDCScopeClient, error) {
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiClient, err := internal.NewClientWithResponses(
		cfg.BaseURL+"/api/server/v1",
		internal.WithHTTPClient(cfg.HTTPClient),
		internal.WithRequestEditorFn(typedAuthEditorFn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC scope client: %w", err)
	}

	return &OIDCScopeClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *OIDCScopeClient) List(ctx context.Context) (*[]OIDCScopeResponseModel, error) {
	resp, err := c.apiClient.GetScopesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get OIDC scopes: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get OIDC scopes: %s", resp.Body)
	}

	return resp.JSON200, nil
}
