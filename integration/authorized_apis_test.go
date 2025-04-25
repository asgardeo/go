//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/asgardeo/go/management"
)

// TestAuthorizedAPIs covers listing, updating, retrieving and deleting authorized APIs for an application.
func TestAuthorizedAPIs(t *testing.T) {
	apiID := os.Getenv("TEST_API_ID")
	if apiID == "" {
		t.Skip("Skipping authorized-apis integration: TEST_API_ID not set")
	}
	scopesEnv := os.Getenv("TEST_API_SCOPES")
	var scopes []string
	if scopesEnv != "" {
		scopes = strings.Split(scopesEnv, ",")
	}

	// create a unique application
	name := fmt.Sprintf("go-asgardeo-auth-%d", time.Now().UnixNano())
	input := management.ApplicationCreateInput{
		Name:                         name,
		AdvancedConfigurations:       &management.AdvancedConfigurations{SkipLoginConsent: true, SkipLogoutConsent: true},
		TemplateID:                   "custom-application-oidc",
		AssociatedRoles:              &management.AssociatedRoles{AllowedAudience: "APPLICATION", Roles: []string{}},
		InboundProtocolConfiguration: &management.InboundProtocolConfiguration{OIDC: &management.InboundOIDCConfig{GrantTypes: []string{"client_credentials"}}},
	}
	app, err := client.Applications().Create(ctx, input)
	if err != nil {
		t.Fatalf("Create app error: %v", err)
	}
	if app.ID == "" {
		t.Fatalf("expected non-empty app ID")
	}
	defer client.Applications().Delete(ctx, app.ID)

	svc := client.AuthorizedAPIs(app.ID)

	// initial list should not contain our apiID
	list0, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List initial failed: %v", err)
	}
	for _, a := range list0 {
		if a.APIID == apiID {
			t.Fatalf("unexpected existing authorized API %s", apiID)
		}
	}

	// update authorized APIs
	toAuth := []management.AuthorizedAPI{{APIID: apiID, Scopes: scopes}}
	if err := svc.Update(ctx, toAuth); err != nil {
		t.Fatalf("Update authorized APIs failed: %v", err)
	}

	// list after update
	list1, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List after update failed: %v", err)
	}
	found := false
	for _, a := range list1 {
		if a.APIID == apiID {
			found = true
			if len(scopes) > 0 && strings.Join(a.Scopes, ",") != strings.Join(scopes, ",") {
				t.Errorf("Scopes mismatch: got %v, want %v", a.Scopes, scopes)
			}
		}
	}
	if !found {
		t.Errorf("authorized API %s not found after update", apiID)
	}

	// get specific authorized API
	got, err := svc.Get(ctx, apiID)
	if err != nil {
		t.Fatalf("Get authorized API failed: %v", err)
	}
	if got.APIID != apiID {
		t.Errorf("Get APIID mismatch: got %s", got.APIID)
	}

	// delete authorized API
	if err := svc.Delete(ctx, apiID); err != nil {
		t.Fatalf("Delete authorized API failed: %v", err)
	}

	// final list should not contain it
	list2, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List after delete failed: %v", err)
	}
	for _, a := range list2 {
		if a.APIID == apiID {
			t.Errorf("authorized API %s still present after delete", apiID)
		}
	}
}
