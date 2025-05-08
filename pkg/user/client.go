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

package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/user/internal"
)

type UserClient struct {
	config    *config.ClientConfig
	apiClient *internal.Client
}

func New(cfg *config.ClientConfig) (*UserClient, error) {
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiClient, err := internal.NewClient(
		cfg.BaseURL+"/scim2",
		internal.WithHTTPClient(cfg.HTTPClient),
		internal.WithRequestEditorFn(typedAuthEditorFn),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user client: %w", err)
	}

	return &UserClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *UserClient) CreateUser(ctx context.Context, user UserCreateModel) (*http.Response, error) {
	creationData := convertToAddUserJSONBodyModel(user)
	resp, err := c.apiClient.AddUser(ctx, creationData)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to create user: %s", resp.Status)
	}
	return resp, nil
}
