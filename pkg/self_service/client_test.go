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

package self_service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *SelfServiceClient {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg := config.DefaultClientConfig().
		WithBaseURL(server.URL).
		WithToken("test-token")
	client, err := New(cfg)
	require.NoError(t, err)
	return client
}

func TestListSessions(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/users/v1/me/sessions", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"userId":"u1","sessions":[{"id":"s1","userAgent":"Firefox","loginTime":"1565586399249"}]}`))
	})

	resp, err := client.ListSessions(context.Background())

	require.NoError(t, err)
	require.Len(t, *resp.Sessions, 1)
	assert.Equal(t, "s1", *(*resp.Sessions)[0].Id)
}

func TestContextToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer user-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	cfg := config.DefaultClientConfig().
		WithBaseURL(server.URL).
		WithContextToken()
	client, err := New(cfg)
	require.NoError(t, err)

	// The request uses the token set on its context
	ctx := config.ContextWithToken(context.Background(), "user-token")
	require.NoError(t, client.TerminateSession(ctx, "s1"))

	// A request without a token on its context fails before it is sent
	require.Error(t, client.TerminateSession(context.Background(), "s1"))
}

func TestContextTokenIgnoredForOtherAuthMethods(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := config.ContextWithToken(context.Background(), "user-token")
	err := client.TerminateSession(ctx, "s1")

	require.NoError(t, err)
}

func TestListAssociationsNoContent(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	associations, err := client.ListAssociations(context.Background())

	require.NoError(t, err)
	assert.Empty(t, associations)
}

func TestErrorWrapsAPIError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":"404","message":"Not found"}`))
	})

	_, err := client.ListConsents(context.Background(), nil)

	var apiErr *common.APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Contains(t, string(apiErr.Body), "Not found")
}
