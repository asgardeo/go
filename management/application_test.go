package management

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplicationService_CRUD(t *testing.T) {
	// Setup a test HTTP server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/server/v1/applications":
			resp := ListApplicationsResponse{
				Count:        1,
				TotalResults: 1,
				Applications: []Application{{ID: "app1", Name: "TestApp", Description: "desc"}},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)

		case r.Method == "GET" && r.URL.Path == "/api/server/v1/applications/app1":
			app := Application{ID: "app1", Name: "TestApp", Description: "desc"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(app)

		case r.Method == "POST" && r.URL.Path == "/api/server/v1/applications":
			var inp ApplicationCreateInput
			json.NewDecoder(r.Body).Decode(&inp)
			respApp := Application{ID: "newId", Name: inp.Name}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respApp)

		case r.Method == "PUT" && r.URL.Path == "/api/server/v1/applications/app1":
			var input Application
			json.NewDecoder(r.Body).Decode(&input)
			input.ID = "app1"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(input)

		case r.Method == "DELETE" && r.URL.Path == "/api/server/v1/applications/app1":
			w.WriteHeader(http.StatusNoContent)

		default:
			t.Fatalf("unexpected method %s or path %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	// Initialize client with static token and custom HTTP client
	client, err := New(
		srv.URL,
		WithStaticToken("token"),
		WithHTTPClient(srv.Client()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	// Test List
	listResp, err := client.Applications().List(ctx, nil)
	if err != nil {
		t.Errorf("List error: %v", err)
	}
	if listResp.TotalResults != 1 || len(listResp.Applications) != 1 || listResp.Applications[0].ID != "app1" {
		t.Errorf("List mismatch: %+v", listResp)
	}

	// Test Get
	gotApp, err := client.Applications().Get(ctx, "app1")
	if err != nil {
		t.Errorf("Get error: %v", err)
	}
	if gotApp.ID != "app1" {
		t.Errorf("Get ID mismatch: got %s", gotApp.ID)
	}

	// Test Create
	input := ApplicationCreateInput{Name: "NewApp"}
	created, err := client.Applications().Create(ctx, input)
	if err != nil {
		t.Errorf("Create error: %v", err)
	}
	if created.ID != "newId" {
		t.Errorf("Create ID mismatch: got %s", created.ID)
	}

	// Test Update
	upd := Application{Name: "UpdatedApp", Description: "updated desc"}
	updated, err := client.Applications().Update(ctx, "app1", upd)
	if err != nil {
		t.Errorf("Update error: %v", err)
	}
	if updated.ID != "app1" || updated.Name != "UpdatedApp" {
		t.Errorf("Update mismatch: %+v", updated)
	}

	// Test Delete
	if err := client.Applications().Delete(ctx, "app1"); err != nil {
		t.Errorf("Delete error: %v", err)
	}
}
