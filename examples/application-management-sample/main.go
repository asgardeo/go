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
		"https://api.asgardeo.io/t/sahan1001",
		management.WithClientCredentials(ctx,
			"bk9jKLdRl9pYquj1hkWZDoqujVIa",
			"6v2OGSM8gRM_jUi1IFOlLNPytf1iAK5uuMcifOJuvsAa",
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
