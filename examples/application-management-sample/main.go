package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shashimalcse/go-asgardeo/management"
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
}

// How to run this example:
// 1. Replace <organization>, <client_id>, and <client_secret> with your actual values.
// 		Ensure the application has the necessary permissions to access the management API.
// 2. Save this code in a file named `main.go`.
// 3. Run the code using `go run main.go`.
