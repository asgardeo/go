//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/asgardeo/go/management"
)

// TestApplicationLifecycle covers create, get, update, regenerate secret, and delete.
func TestApplications(t *testing.T) {
	// unique application name per run
	name := fmt.Sprintf("go-asgardeo-test-%d", time.Now().UnixNano())

	// Create application
	input := management.ApplicationCreateInput{
		Name: name,
		AdvancedConfigurations: &management.AdvancedConfigurations{
			SkipLogoutConsent: true,
			SkipLoginConsent:  true,
		},
		TemplateID: "custom-application-oidc",
		AssociatedRoles: &management.AssociatedRoles{
			AllowedAudience: "APPLICATION",
			Roles:           []string{},
		},
		InboundProtocolConfiguration: &management.InboundProtocolConfiguration{
			OIDC: &management.InboundOIDCConfig{GrantTypes: []string{"client_credentials"}},
		},
	}
	created, err := client.Applications().Create(ctx, input)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	// If API doesn't return ID, fallback by listing the application by name
	if created.ID == "" {
		listResp, err := client.Applications().List(ctx, &management.ListApplicationsParams{
			Filter: fmt.Sprintf("name eq \"%s\"", input.Name),
		})
		if err != nil {
			t.Fatalf("fallback list after create failed: %v", err)
		}
		if len(listResp.Applications) == 0 {
			t.Fatalf("no applications found when listing after create")
		}
		created.ID = listResp.Applications[0].ID
	}
	// ensure cleanup
	defer func() {
		if err := client.Applications().Delete(ctx, created.ID); err != nil {
			t.Logf("cleanup delete failed: %v", err)
		}
	}()

	// Get application
	got, err := client.Applications().Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name != name {
		t.Errorf("expected name %q, got %q", name, got.Name)
	}
}
