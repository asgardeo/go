/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package sdk

import (
	"github.com/asgardeo/go/pkg/api_resource"
	"github.com/asgardeo/go/pkg/application"
	"github.com/asgardeo/go/pkg/authenticator"
	"github.com/asgardeo/go/pkg/claim"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/identity_provider"
)

// Client is the main SDK client that provides access to all service clients
type Client struct {
	Config           *config.ClientConfig
	Application      *application.ApplicationClient
	APIResource      *api_resource.APIResourceClient
	IdentityProvider *identity_provider.IdentityProviderClient
	Authenticator    *authenticator.AuthenticatorClient
	Claim            *claim.ClaimClient
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

	identityProviderClient, err := identity_provider.New(cfg)
	if err != nil {
		return nil, err
	}

	authenticatorClient, err := authenticator.New(cfg)
	if err != nil {
		return nil, err
	}

	claimClient, err := claim.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:           cfg,
		Application:      appClient,
		APIResource:      apiResourceClient,
		IdentityProvider: identityProviderClient,
		Authenticator:    authenticatorClient,
		Claim:            claimClient,
	}, nil
}
