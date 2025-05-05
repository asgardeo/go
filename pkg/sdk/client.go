package sdk

import (
	"github.com/asgardeo/go/pkg/api_resource"
	"github.com/asgardeo/go/pkg/application"
	"github.com/asgardeo/go/pkg/config"
)

// Client is the main SDK client that provides access to all service clients
type Client struct {
	Config      *config.ClientConfig
	Application *application.ApplicationClient
	APIResource *api_resource.APIResourceClient
}

// NewClient creates a new SDK client with the given configuration
func New(cfg *config.ClientConfig) (*Client, error) {

	appClient, err := application.New(cfg)
	if err != nil {
		return nil, err
	}

	apiResourceClient, err := api_resource.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:      cfg,
		Application: appClient,
		APIResource: apiResourceClient,
	}, nil
}
