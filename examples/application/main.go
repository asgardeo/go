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

	// Create a configuration with a client credentials grant type
	cfg := config.DefaultClientConfig().
		WithBaseURL("https://api.asgardeo.io/t/<tenant-domain>").
		WithTimeout(10*time.Second).
		WithClientCredentials(
			"client_id",
			"client_secret",
		)

	// Create a client with the client credentials configuration
	client, err := sdk.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Use the client with token authentication
	ctx := context.Background()
	apps, err := client.ApplicationClient.List(ctx, 10, 0)
	if err != nil {
		log.Printf("Error listing users: %v", err)
	} else {
		fmt.Printf("Found %d applications\n", len(*apps.Applications))
	}
}
