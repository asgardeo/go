package common

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/config"
)

// CreateAuthRequestEditorFunc returns a function that adds authentication to requests
func CreateAuthRequestEditorFunc(cfg *config.ClientConfig) interface{} {
	// Return a function that matches the RequestEditorFn signature
	// The caller will need to cast this to the appropriate type
	return func(ctx context.Context, req *http.Request) error {

		token, err := cfg.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get authentication token: %w", err)
		}

		if token != "" {
			// Add Authorization header with Bearer token
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		}
		return nil
	}
}
