package management

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper for string pointers
func ptr(s string) *string { return &s }

func TestAPIResourceService_CRUD(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/server/v1/api-resources":
			resp := struct {
				Resources []APIResource `json:"resources"`
			}{
				Resources: []APIResource{{
					ID: "res1", Name: "TestRes", Identifier: "testRes", Description: "desc",
					RequiresAuthorization: false,
				}},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)

		case r.Method == "GET" && r.URL.Path == "/api/server/v1/api-resources/res1":
			res := APIResource{
				ID: "res1", Name: "TestRes", Identifier: "testRes", Description: "desc",
				RequiresAuthorization: false,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)

		case r.Method == "POST" && r.URL.Path == "/api/server/v1/api-resources":
			var inp APIResourceCreateInput
			json.NewDecoder(r.Body).Decode(&inp)
			if inp.Identifier != "greetings_api" || !inp.RequiresAuthorization {
				t.Errorf("invalid create input: %+v", inp)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(APIResource{ID: "newId", Name: inp.Name})

		case r.Method == "PUT" && r.URL.Path == "/api/server/v1/api-resources/res1":
			var inp APIResource
			json.NewDecoder(r.Body).Decode(&inp)
			inp.ID = "res1"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(inp)

		case r.Method == "DELETE" && r.URL.Path == "/api/server/v1/api-resources/res1":
			w.WriteHeader(http.StatusNoContent)

		case r.Method == "PATCH" && r.URL.Path == "/api/server/v1/api-resources/res1":
			// accept patch request, respond with no content
			w.WriteHeader(http.StatusNoContent)

		default:
			t.Fatalf("unexpected method %s or path %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client, err := New(srv.URL, WithStaticToken("token"), WithHTTPClient(srv.Client()))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	// List
	list, err := client.APIResources().List(ctx)
	if err != nil {
		t.Errorf("List error: %v", err)
	}
	if len(list) != 1 || list[0].ID != "res1" {
		t.Errorf("List mismatch: %+v", list)
	}

	// Get
	got, err := client.APIResources().Get(ctx, "res1")
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if got.ID != "res1" {
		t.Errorf("Get ID mismatch: got %s", got.ID)
	}

	// Create
	created, err := client.APIResources().Create(ctx, APIResourceCreateInput{
		Identifier: "greetings_api", Name: "NewRes", RequiresAuthorization: true,
	})
	if err != nil {
		t.Errorf("Create error: %v", err)
	}
	if created.ID != "newId" {
		t.Errorf("Create ID mismatch: got %s", created.ID)
	}

	// Update
	newName := "UpdatedRes"
	updInput := APIResourceUpdateInput{Name: &newName, Scopes: []string{"write"}}
	updated, err := client.APIResources().Update(ctx, "res1", updInput)
	if err != nil {
		t.Errorf("Update error: %v", err)
	}
	if updated.ID != "res1" || updated.Name != "UpdatedRes" {
		t.Errorf("Update mismatch: %+v", updated)
	}

	// Delete
	if err := client.APIResources().Delete(ctx, "res1"); err != nil {
		t.Errorf("Delete error: %v", err)
	}

	// Patch
	if err := client.APIResources().Patch(ctx, "res1", APIResourcePatchModel{
		Name: ptr("PatchedName"), Description: ptr("patched desc"),
		AddedScopes: []ScopeCreationModel{{Name: "new:scope", DisplayName: "New Scope", Description: "desc"}},
	}); err != nil {
		t.Errorf("Patch error: %v", err)
	}
}

func TestAPIResourceScopes_CRUD(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/server/v1/api-resources/res1/scopes":
			// return existing scopes
			sc := []ScopeGetModel{{ID: "s1", Name: "read", DisplayName: "Read", Description: "desc"}}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sc)

		case r.Method == "PUT" && r.URL.Path == "/api/server/v1/api-resources/res1/scopes":
			var inp []ScopeCreationModel
			json.NewDecoder(r.Body).Decode(&inp)
			// accept and return no content
			w.WriteHeader(http.StatusNoContent)

		case r.Method == "DELETE" && r.URL.Path == "/api/server/v1/api-resources/res1/scopes/read":
			w.WriteHeader(http.StatusNoContent)

		case r.Method == "GET" && r.URL.Path == "/scopes":
			// global scopes list
			scList := []ScopeGetModel{{ID: "s2", Name: "admin", DisplayName: "Admin"}}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(scList)

		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client, err := New(srv.URL, WithStaticToken("token"), WithHTTPClient(srv.Client()))
	if err != nil {
		t.Fatalf("client init failed: %v", err)
	}
	ctx := context.Background()

	// GetScopes
	scopes, err := client.APIResources().GetScopes(ctx, "res1")
	if err != nil {
		t.Errorf("GetScopes error: %v", err)
	}
	if len(scopes) != 1 || scopes[0].Name != "read" || scopes[0].ID != "s1" {
		t.Errorf("GetScopes mismatch: %+v", scopes)
	}

	// AddScopes
	err = client.APIResources().AddScopes(ctx, "res1", []ScopeCreationModel{{Name: "write", DisplayName: "Write"}})
	if err != nil {
		t.Errorf("AddScopes error: %v", err)
	}

	// DeleteScope
	if err := client.APIResources().DeleteScope(ctx, "res1", "read"); err != nil {
		t.Errorf("DeleteScope error: %v", err)
	}

	// ListScopes
	sAll, err := client.APIResources().ListScopes(ctx, "")
	if err != nil {
		t.Errorf("ListScopes error: %v", err)
	}
	if len(sAll) != 1 || sAll[0].Name != "admin" || sAll[0].ID != "s2" {
		t.Errorf("ListScopes mismatch: %+v", sAll)
	}
}
