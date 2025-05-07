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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asgardeo/go/pkg/application/internal"
	"github.com/asgardeo/go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApplicationClient(t *testing.T) {

	// Test creating a new client with custom config
	customCfg := config.DefaultClientConfig().WithBaseURL("https://api.asgardeo.io/t/test-domain")
	client, err := New(customCfg)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, customCfg, client.config)
}

func TestListApplications(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/api/server/v1/applications", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "10", query.Get("limit"))
		assert.Equal(t, "0", query.Get("offset"))

		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer test-token", authHeader)

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create mock data
		app1 := internal.ApplicationListItem{
			Id:       stringPtr("app-id-1"),
			Name:     stringPtr("Test App 1"),
			ClientId: stringPtr("client-1"),
		}
		app2 := internal.ApplicationListItem{
			Id:       stringPtr("app-id-2"),
			Name:     stringPtr("Test App 2"),
			ClientId: stringPtr("client-2"),
		}

		apps := []internal.ApplicationListItem{app1, app2}
		response := internal.ApplicationListResponse{
			Applications: &apps,
			Count:        intPtr(2),
			TotalResults: intPtr(2),
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test server URL
	cfg := config.DefaultClientConfig().
		WithBaseURL(server.URL).
		WithToken("test-token")

	client, err := New(cfg)
	require.NoError(t, err)

	// Call method to test
	limit := 10
	offset := 0
	resp, err := client.List(context.Background(), limit, offset)

	// Assert expectations
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, *resp.Applications, 2)
}
