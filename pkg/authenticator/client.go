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

package authenticator

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

type AuthenticatorClient struct {
	config    *config.ClientConfig
	apiClient *ClientWithResponses
}

func New(cfg *config.ClientConfig) (*AuthenticatorClient, error) {
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
		return nil, fmt.Errorf("Failed to create authenticator client: %w", err)
	}

	return &AuthenticatorClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *AuthenticatorClient) List(ctx context.Context, params *GetAuthenticatorsParams) (*Authenticators, error) {
	resp, err := c.apiClient.GetAuthenticatorsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Failed to list authenticators: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list authenticators: %s", resp.Body)
	}

	return resp.JSON200, nil
}

func (c *AuthenticatorClient) ListLocalAuthenticators(ctx context.Context) (*Authenticators, error) {
	resp, err := c.apiClient.GetAuthenticatorsWithResponse(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to list local authenticators: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list local authenticators: %s", resp.Body)
	}
	allAuthenticators := resp.JSON200
	if allAuthenticators == nil {
		return nil, fmt.Errorf("Failed to list local authenticators: %s", resp.Body)
	}
	localAuthenticators := make([]Authenticator, 0)
	for _, authenticator := range *allAuthenticators {
		if *authenticator.Type == "LOCAL" {
			localAuthenticators = append(localAuthenticators, authenticator)
		}
	}
	return &localAuthenticators, nil
}
