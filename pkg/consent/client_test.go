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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asgardeo/go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/identity/consent-mgt/v1.0/consents", r.URL.Path)
		assert.Equal(t, "anne", r.URL.Query().Get("piiPrincipalId"))
		assert.Equal(t, "ACTIVE", r.URL.Query().Get("state"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"consentReceiptID":"r1","piiPrincipalId":"anne","tenantDomain":"carbon.super","state":"ACTIVE","spDisplayName":"Pickup"}]`))
	}))
	defer server.Close()

	client, err := New(config.DefaultClientConfig().WithBaseURL(server.URL).WithToken("test-token"))
	require.NoError(t, err)

	principal, state := "anne", "ACTIVE"
	consents, err := client.List(context.Background(), &ConsentListParamsModel{PiiPrincipalId: &principal, State: &state})

	require.NoError(t, err)
	require.Len(t, consents, 1)
	assert.Equal(t, "r1", consents[0].ConsentReceiptID)
	assert.Equal(t, "Pickup", *consents[0].SpDisplayName)
}

func TestRevokeReceipt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/identity/consent-mgt/v1.0/consents/receipts/r1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New(config.DefaultClientConfig().WithBaseURL(server.URL).WithToken("test-token"))
	require.NoError(t, err)

	require.NoError(t, client.RevokeReceipt(context.Background(), "r1"))
}
