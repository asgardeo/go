package management

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrganizationService_List(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

		if r.URL.Path != "/api/server/v1/organizations" {
			t.Errorf("Expected path '/api/server/v1/organizations', got '%s'", r.URL.Path)
		}

		// Check query parameters when provided
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected query param limit=10, got '%s'", r.URL.Query().Get("limit"))
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"count": 2,
			"totalResults": 2,
			"organizations": [
				{
					"id": "org-123",
					"name": "test-org-1",
					"displayName": "Test Organization 1",
					"description": "This is test organization 1"
				},
				{
					"id": "org-456",
					"name": "test-org-2",
					"displayName": "Test Organization 2",
					"description": "This is test organization 2"
				}
			]
		}`))
	}))
	defer ts.Close()

	// Create client with test server URL
	client, err := New(ts.URL, WithStaticToken("test-token"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the method with parameters
	params := &ListOrganizationsParams{
		Limit:     10,
		Recursive: true,
	}

	resp, err := client.Organizations().List(context.Background(), params)
	if err != nil {
		t.Fatalf("Organizations.List returned error: %v", err)
	}

	// Verify response
	if resp.Count != 2 {
		t.Errorf("Expected count=2, got %d", resp.Count)
	}

	if len(resp.Organizations) != 2 {
		t.Fatalf("Expected 2 organizations, got %d", len(resp.Organizations))
	}

	org := resp.Organizations[0]
	if org.ID != "org-123" {
		t.Errorf("Expected org ID org-123, got %s", org.ID)
	}

	if org.Name != "test-org-1" {
		t.Errorf("Expected org name test-org-1, got %s", org.Name)
	}
}

func TestOrganizationService_Get(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

		if r.URL.Path != "/api/server/v1/organizations/org-123" {
			t.Errorf("Expected path '/api/server/v1/organizations/org-123', got '%s'", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "org-123",
			"name": "test-org",
			"displayName": "Test Organization",
			"description": "This is a test organization"
		}`))
	}))
	defer ts.Close()

	// Create client with test server URL
	client, err := New(ts.URL, WithStaticToken("test-token"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the method
	org, err := client.Organizations().Get(context.Background(), "org-123")
	if err != nil {
		t.Fatalf("Organizations.Get returned error: %v", err)
	}

	// Verify response
	if org.ID != "org-123" {
		t.Errorf("Expected org ID org-123, got %s", org.ID)
	}

	if org.Name != "test-org" {
		t.Errorf("Expected org name test-org, got %s", org.Name)
	}

	if org.DisplayName != "Test Organization" {
		t.Errorf("Expected org display name 'Test Organization', got '%s'", org.DisplayName)
	}

	if org.Description != "This is a test organization" {
		t.Errorf("Expected org description 'This is a test organization', got '%s'", org.Description)
	}
}

func TestOrganizationService_Create(t *testing.T) {
	// Setup test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		if r.URL.Path != "/api/server/v1/organizations" {
			t.Errorf("Expected path '/api/server/v1/organizations', got '%s'", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": "org-new",
			"name": "new-org",
			"displayName": "New Organization",
			"description": "This is a new test organization"
		}`))
	}))
	defer ts.Close()

	// Create client with test server URL
	client, err := New(ts.URL, WithStaticToken("test-token"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the method
	input := OrganizationCreateInput{
		Name:        "new-org",
		DisplayName: "New Organization",
		Description: "This is a new test organization",
	}

	org, err := client.Organizations().Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Organizations.Create returned error: %v", err)
	}

	// Verify response
	if org.ID != "org-new" {
		t.Errorf("Expected org ID org-new, got %s", org.ID)
	}

	if org.Name != "new-org" {
		t.Errorf("Expected org name new-org, got %s", org.Name)
	}

	if org.DisplayName != "New Organization" {
		t.Errorf("Expected org display name 'New Organization', got '%s'", org.DisplayName)
	}
}

func TestOrganizationService_ErrorHandling(t *testing.T) {
	// Setup test server that returns an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid_request"}`))
	}))
	defer ts.Close()

	// Create client with test server URL
	client, err := New(ts.URL, WithStaticToken("test-token"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test List with error
	_, err = client.Organizations().List(context.Background(), nil)
	if err == nil {
		t.Error("Expected error from Organizations.List, got nil")
	}

	// Test Get with error
	_, err = client.Organizations().Get(context.Background(), "org-123")
	if err == nil {
		t.Error("Expected error from Organizations.Get, got nil")
	}

	// Test Create with error
	input := OrganizationCreateInput{Name: "test-org"}
	_, err = client.Organizations().Create(context.Background(), input)
	if err == nil {
		t.Error("Expected error from Organizations.Create, got nil")
	}
}
