package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/sdk"
)

func main() {

	// Initialize the client configurations.
	cfg := config.DefaultClientConfig().
		WithBaseURL("https://api.asgardeo.io/t/<tenant-domain>").
		WithTimeout(10*time.Second).
		WithClientCredentials(
			"client_id",
			"client_secret",
		)

	// Create a client with the given configurations.
	client, err := sdk.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Use the client with token authentication.
	ctx := context.Background()
	apps, err := client.APIResource.List(ctx, nil, nil, nil, nil)
	if err != nil {
		log.Printf("Error listing API Resources: %v", err)
	} else {
		fmt.Printf("Found %d API Resources\n", len(*apps.APIResources))
	}
}
