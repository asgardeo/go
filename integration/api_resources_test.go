//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/shashimalcse/go-asgardeo/management"
)

// TestAPIResourceLifecycle performs full CRUD and scope operations against real Asgardeo server
func TestAPIResources(t *testing.T) {
	unique := fmt.Sprintf("go-asgardeo-test-%d", time.Now().UnixNano())
	// Create
	createInput := management.APIResourceCreateInput{
		Identifier:            unique,
		Name:                  unique,
		Description:           "Integration test resource",
		RequiresAuthorization: false,
		Scopes:                []management.ScopeCreationModel{},
	}
	res, err := client.APIResources().Create(ctx, createInput)
	if err != nil {
		t.Fatalf("Create API resource failed: %v", err)
	}
	if res.ID == "" {
		t.Fatalf("Create returned empty ID")
	}
	// ensure cleanup
	defer func() {
		if err := client.APIResources().Delete(ctx, res.ID); err != nil {
			t.Logf("cleanup delete failed: %v", err)
		}
	}()

	// Get
	got, err := client.APIResources().Get(ctx, res.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name != unique {
		t.Errorf("expected name %q, got %q", unique, got.Name)
	}

	// Patch: update description and add a scope
	newDesc := "patched desc"
	newScope := management.ScopeCreationModel{Name: unique + ":scope", DisplayName: "Test Scope"}
	if err := client.APIResources().Patch(ctx, res.ID, management.APIResourcePatchModel{
		Description: &newDesc,
		AddedScopes: []management.ScopeCreationModel{newScope},
	}); err != nil {
		t.Fatalf("Patch failed: %v", err)
	}

	// Verify patch applied
	got2, err := client.APIResources().Get(ctx, res.ID)
	if err != nil {
		t.Fatalf("Get after patch failed: %v", err)
	}
	if got2.Description != newDesc {
		t.Errorf("patched description mismatch: %q", got2.Description)
	}

	// GetScopes
	scopes, err := client.APIResources().GetScopes(ctx, res.ID)
	if err != nil {
		t.Fatalf("GetScopes failed: %v", err)
	}
	found := false
	for _, s := range scopes {
		if s.Name == newScope.Name {
			found = true
		}
	}
	if !found {
		t.Errorf("Scope %q not found in API resource scopes", newScope.Name)
	}

	// DeleteScope
	if err := client.APIResources().DeleteScope(ctx, res.ID, newScope.Name); err != nil {
		t.Errorf("DeleteScope failed: %v", err)
	}

	// AddScopes
	addScopes := []management.ScopeCreationModel{
		{Name: unique + ":s1", DisplayName: "S1"},
		{Name: unique + ":s2", DisplayName: "S2"},
	}
	if err := client.APIResources().AddScopes(ctx, res.ID, addScopes); err != nil {
		t.Errorf("AddScopes failed: %v", err)
	}

	sc2, err := client.APIResources().GetScopes(ctx, res.ID)
	if err != nil {
		t.Fatalf("GetScopes after add failed: %v", err)
	}
	if len(sc2) != len(addScopes) {
		t.Errorf("expected %d scopes, got %d", len(addScopes), len(sc2))
	}

	// ListScopes (global)
	gl, err := client.APIResources().ListScopes(ctx, "")
	if err != nil {
		t.Errorf("ListScopes (global) failed: %v", err)
	}
	if len(gl) == 0 {
		t.Errorf("expected global scopes to be non-empty")
	}
}
