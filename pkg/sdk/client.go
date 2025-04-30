package sdk

import (
	"github.com/asgardeo/go/pkg/application"
	"github.com/asgardeo/go/pkg/config"
)

// Client is the main SDK client that provides access to all service clients
type Client struct {
	Config            *config.ClientConfig
	ApplicationClient *application.ApplicationClient
}

// NewClient creates a new SDK client with the given configuration
func NewClient(cfg *config.ClientConfig) (*Client, error) {

	appClient, err := application.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:            cfg,
		ApplicationClient: appClient,
	}, nil
}
