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

package claim

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/claim/internal"
	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

type ClaimClient struct {
	config    *config.ClientConfig
	apiClient *internal.ClientWithResponses
}

func New(cfg *config.ClientConfig) (*ClaimClient, error) {
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
		return nil, fmt.Errorf("failed to create claim client: %w", err)
	}

	return &ClaimClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *ClaimClient) ListLocalClaims(ctx context.Context, params *LocalClaimListParamsModel) (*[]LocalClaimResponseModel, error) {
	resp, err := c.apiClient.GetLocalClaimsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %w", err)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to get claims: %s", resp.Status())
	}

	return resp.JSON200, nil
}

func (c *ClaimClient) ListExternalClaims(ctx context.Context, dialectId string, params *ExternalClaimListParamsModel) (*[]ExternalClaimResponseModel, error) {
	resp, err := c.apiClient.GetExternalClaimsWithResponse(ctx, dialectId, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get external claims: %w", err)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to get external claims: %s", resp.Status())
	}

	return resp.JSON200, nil
}

func (c *ClaimClient) ListOIDCClaims(ctx context.Context, params *ExternalClaimListParamsModel) (*[]ExternalClaimResponseModel, error) {
	resp, err := c.apiClient.GetExternalClaimsWithResponse(ctx, ClaimDialectIDs.OIDC, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get OIDC claims: %w", err)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to get OIDC claims: %s", resp.Status())
	}

	return resp.JSON200, nil
}
