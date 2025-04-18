//go:build integration
// +build integration

package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/thilinashashimalsenarath/go-asgardeo/management"
)

var (
	client *management.Client
	ctx    context.Context
)

func TestMain(m *testing.M) {
	// Load .env if present in cwd
	_ = godotenv.Load()
	// Then try parent directory for .env
	if err := godotenv.Load("../.env"); err == nil {
		log.Println("Loaded .env from parent directory")
	}
	baseURL := os.Getenv("ASGARDEO_BASE_URL")
	clientID := os.Getenv("ASGARDEO_CLIENT_ID")
	clientSecret := os.Getenv("ASGARDEO_CLIENT_SECRET")
	log.Printf("Env loaded: ASGARDEO_BASE_URL=%q, CLIENT_ID=%q, CLIENT_SECRET=%q", baseURL, clientID, clientSecret)
	if baseURL == "" || clientID == "" || clientSecret == "" {
		log.Println("Skipping integration tests: missing ASGARDEO env vars")
		os.Exit(0)
	}
	ctx = context.Background()
	var err error
	client, err = management.New(
		baseURL,
		management.WithClientCredentials(ctx, clientID, clientSecret),
	)
	if err != nil {
		log.Fatalf("failed to initialize management client: %v", err)
	}
	os.Exit(m.Run())
}
