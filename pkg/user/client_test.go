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

package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asgardeo/go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *UserClient {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg := config.DefaultClientConfig().
		WithBaseURL(server.URL).
		WithToken("test-token")
	client, err := New(cfg)
	require.NoError(t, err)
	return client
}

func TestCreateUser(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/scim2/Users", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]interface{}{
			"userName": "DEFAULT/anne@x.com",
			"email":    "anne@x.com",
			"password": "Secret#123",
			"name":     map[string]interface{}{"givenName": "Anne", "familyName": "Perera"},
			"emails":   []interface{}{map[string]interface{}{"primary": true, "value": "anne@x.com"}},
		}, body)

		w.WriteHeader(http.StatusCreated)
	})

	_, err := client.CreateUser(context.Background(), UserCreateModel{
		Username:  "DEFAULT/anne@x.com",
		Email:     "anne@x.com",
		Password:  "Secret#123",
		FirstName: "Anne",
		LastName:  "Perera",
	})

	require.NoError(t, err)
}

func TestGetMe(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/scim2/Me", r.URL.Path)
		assert.Equal(t, "userName,name", r.URL.Query().Get("attributes"))

		w.Header().Set("Content-Type", "application/scim+json")
		w.Write([]byte(`{"userName":"anne","name":{"givenName":"Anne"}}`))
	})

	me, err := client.GetMe(context.Background(), "userName,name")

	require.NoError(t, err)
	assert.Equal(t, "anne", me["userName"])
}

func TestPatchMe(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/scim2/Me", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "application/scim+json", r.Header.Get("Content-Type"))

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []interface{}{scimPatchOpSchema}, body["schemas"])
		assert.Equal(t, []interface{}{map[string]interface{}{
			"op": "replace", "value": map[string]interface{}{"name": map[string]interface{}{"givenName": "Anna"}},
		}}, body["Operations"])

		w.Header().Set("Content-Type", "application/scim+json")
		w.Write([]byte(`{"userName":"anne","name":{"givenName":"Anna"}}`))
	})

	value := map[string]interface{}{"name": map[string]interface{}{"givenName": "Anna"}}
	me, err := client.PatchMe(context.Background(), []PatchOperationModel{
		{Op: PatchOpReplace, Value: value},
	})

	require.NoError(t, err)
	assert.Equal(t, "Anna", me["name"].(map[string]interface{})["givenName"])
}
