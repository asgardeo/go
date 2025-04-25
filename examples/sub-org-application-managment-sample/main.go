package main

import (
	"context"
	"fmt"
	"log"

	"github.com/asgardeo/go/management"
)

func main() {
	// Initialize management client with client credentials
	ctx := context.Background()
	client, err := management.New(
		"https://api.asgardeo.io/t/<organization>",
		management.WithClientCredentials(ctx,
			"<client_id>",
			"<client_secret>",
		),
	)

	if err != nil {
		log.Fatalf("Failed to create management client: %v", err)
	}
	// List all applications
	apps, err := client.Applications().List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list applications: %v", err)
	}
	for _, app := range apps.Applications {
		fmt.Printf("Application ID: %s, Name: %s\n", app.ID, app.Name)
	}

	// --- Sample: List applications in a sub-organization ---
	// Replace <sub_org_id> with your sub-organization ID.
	subOrgID := "<sub_org_id>"

	// Create a SubOrgApplicationService using the management client.
	subOrgAppService := client.SubOrgApplications()

	// Optionally, provide query parameters (e.g., limit, offset, filter).
	subOrgApps, err := subOrgAppService.GetAll(subOrgID, nil)
	if err != nil {
		log.Fatalf("Failed to list sub-org applications: %v", err)
	}
	for _, app := range subOrgApps.Applications {
		fmt.Printf("[SubOrg] Application ID: %s, Name: %s\n", app.ID, app.Name)
	}
}

// How to run this example:
// 1. Replace <organization>, <client_id>, and <client_secret> and <sub_org_id> with your actual values.
// 		Ensure the application has the necessary permissions to access the management API.
// 2. Save this code in a file named `main.go`.
// 3. Run the code using `go run main.go`.
