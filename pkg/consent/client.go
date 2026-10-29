/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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

package consent

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/consent/internal"
)

// ConsentClient is a wrapper around the generated client for the Consent Management API
type ConsentClient struct {
	config    *config.ClientConfig
	apiClient *internal.ClientWithResponses
}

// New creates a new Consent Management API client
func New(cfg *config.ClientConfig) (*ConsentClient, error) {
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiClient, err := internal.NewClientWithResponses(
		cfg.BaseURL+"/api/identity/consent-mgt/v1.0",
		internal.WithHTTPClient(cfg.HTTPClient),
		internal.WithRequestEditorFn(typedAuthEditorFn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consent client: %w", err)
	}

	return &ConsentClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

// List retrieves consent receipts matching the given filters
func (c *ConsentClient) List(ctx context.Context, params *ConsentListParamsModel) ([]ConsentListItemModel, error) {
	resp, err := c.apiClient.GetConsentsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list consents: %w", err)
	}
	if err := common.CheckStatus("list consents", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []ConsentListItemModel{}, nil
	}
	return *resp.JSON200, nil
}

// GetReceipt retrieves a consent receipt
func (c *ConsentClient) GetReceipt(ctx context.Context, receiptID string) (*ConsentReceiptModel, error) {
	resp, err := c.apiClient.GetConsentsReceiptsReceiptIdWithResponse(ctx, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consent receipt: %w", err)
	}
	if err := common.CheckStatus("get consent receipt", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// RevokeReceipt revokes a consent receipt
func (c *ConsentClient) RevokeReceipt(ctx context.Context, receiptID string) error {
	resp, err := c.apiClient.DeleteConsentsReceiptsReceiptIdWithResponse(ctx, receiptID)
	if err != nil {
		return fmt.Errorf("failed to revoke consent receipt: %w", err)
	}
	return common.CheckStatus("revoke consent receipt", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}
