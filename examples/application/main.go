package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asgardeo/go/pkg/application"
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
	client, err := sdk.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Use the client with token authentication
	ctx := context.Background()

	// List applications.
	apps, err := client.Application.List(ctx, 10, 0)
	if err != nil {
		log.Printf("Error listing users: %v", err)
	} else {
		fmt.Printf("Found %d applications\n", len(*apps.Applications))
	}

	// Authorize API.
	id := "1f616716-f518-48a5-a497-5eb0e2200b4f"
	policyIdentifier := "RBAC"
	scopes := []string{"internal_user_mgt_view", "internal_user_mgt_list"}
	authorizedAPI := application.AddAuthorizedAPIJSONRequestBody{
		Id:               &id,
		PolicyIdentifier: &policyIdentifier,
		Scopes:           &scopes,
	}
	_, err = client.Application.AuthorizeAPI(ctx, "app_uuid", authorizedAPI)
	if err != nil {
		log.Printf("Error authorizing API: %v", err)
	} else {
		log.Printf("API authorized successfully.")
	}
}
