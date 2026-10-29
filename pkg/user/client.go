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
	"encoding/json"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("failed to create user client: %w", err)
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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create user: %s", resp.Status)
	}
	return resp, nil
}

const scimPatchOpSchema = "urn:ietf:params:scim:api:messages:2.0:PatchOp"

// GetMe retrieves the authenticated user, optionally limited to comma-separated attributes
func (c *UserClient) GetMe(ctx context.Context, attributes string) (MeModel, error) {
	resp, err := c.apiClient.GetMe(ctx, attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to get the authenticated user: %w", err)
	}
	me := MeModel{}
	if err := decodeResponse(resp, "get the authenticated user", &me); err != nil {
		return nil, err
	}
	return me, nil
}

// PatchMe updates the authenticated user
func (c *UserClient) PatchMe(ctx context.Context, operations []PatchOperationModel) (MeModel, error) {
	body := internal.PatchMeJSONBody{
		Schemas:    []string{scimPatchOpSchema},
		Operations: operations,
	}
	resp, err := c.apiClient.PatchMe(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("failed to update the authenticated user: %w", err)
	}
	me := MeModel{}
	if err := decodeResponse(resp, "update the authenticated user", &me); err != nil {
		return nil, err
	}
	return me, nil
}

// GetSchemas retrieves the SCIM schemas
func (c *UserClient) GetSchemas(ctx context.Context) (interface{}, error) {
	resp, err := c.apiClient.GetSchemas(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas: %w", err)
	}
	var schemas interface{}
	if err := decodeResponse(resp, "get schemas", &schemas); err != nil {
		return nil, err
	}
	return schemas, nil
}

// decodeResponse decodes a 200 OK JSON response into out
func decodeResponse(resp *http.Response, action string, out interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to %s: %w", action, err)
	}
	if err := common.CheckStatus(action, resp.StatusCode, body, http.StatusOK); err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("failed to %s: %w", action, err)
	}
	return nil
}
