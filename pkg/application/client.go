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

package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
)

const (
	defaultLimit  int = 10
	defaultOffset int = 0
)

// ApplicationClient is a wrapper around the generated client for the Application Management API
type ApplicationClient struct {
	config    *config.ClientConfig
	apiClient *ClientWithResponses
}

type ApplicationBaseModel struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ClientId         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	RedirectURL      string `json:"redirect_url"`
	AuthorizedScopes string `json:"scope"`
}

// New creates a new Application Management API client
func New(cfg *config.ClientConfig) (*ApplicationClient, error) {
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiClient, err := NewClientWithResponses(
		cfg.BaseURL+"/api/server/v1",
		WithHTTPClient(cfg.HTTPClient),
		WithRequestEditorFn(typedAuthEditorFn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create application client: %w", err)
	}

	return &ApplicationClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

func (c *ApplicationClient) List(ctx context.Context, limit, offset int) (*ApplicationListResponse, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = defaultOffset
	}

	params := GetAllApplicationsParams{
		Limit:  &limit,
		Offset: &offset,
	}
	resp, err := c.apiClient.GetAllApplicationsWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("Failed to list applications: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to list applications: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// CreateSinglePageApp creates a new Single Page Application with sensible defaults
func (c *ApplicationClient) CreateSinglePageApp(ctx context.Context, name string, redirectURL string) (*ApplicationBaseModel, error) {
	appRequest, err := c.buildSPARequest(name, redirectURL)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create SPA: %w", err)
	}

	return c.processCreatePublicClientAppResponse(ctx, resp, name, redirectURL)
}

// CreateMobileApp creates a new Mobile Application with sensible defaults
func (c *ApplicationClient) CreateMobileApp(ctx context.Context, name string, redirectURL string) (*ApplicationBaseModel, error) {
	appRequest, err := c.buildMobileAppRequest(name, redirectURL)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create mobile app: %w", err)
	}

	return c.processCreatePublicClientAppResponse(ctx, resp, name, redirectURL)
}

// CreateM2MApp creates a new Machine-to-Machine (M2M) Application
func (c *ApplicationClient) CreateM2MApp(ctx context.Context, name string) (*ApplicationBaseModel, error) {
	appRequest, err := c.buildM2MAppRequest(name)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create M2M app: %w", err)
	}

	return c.processCreateConfidentialClientAppResponse(ctx, resp, name)
}

// GetByName finds an application by name and returns its details
// todo: improve application details being fetched beyond appId, name, clientId and clientSecret
func (c *ApplicationClient) GetByName(ctx context.Context, name string) (*ApplicationBaseModel, error) {
	filter := fmt.Sprintf("name eq %s", name)
	excludeSystemPortals := true

	params := GetAllApplicationsParams{
		Filter:               &filter,
		ExcludeSystemPortals: &excludeSystemPortals,
		Attributes:           stringPtr("templateId,clientId"),
	}

	resp, err := c.apiClient.GetAllApplicationsWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to find application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to find application: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil || resp.JSON200.Applications == nil || len(*resp.JSON200.Applications) == 0 {
		return nil, fmt.Errorf("application with name '%s' not found", name)
	}

	var targetApp *ApplicationListItem
	for _, app := range *resp.JSON200.Applications {
		if app.Name != nil && *app.Name == name {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		return nil, fmt.Errorf("application with name '%s' not found", name)
	}

	return c.getApplicationDetails(ctx, *targetApp.Id)
}

// GetByClienId finds an application by clientId and returns its details
// todo: improve application details being fetched beyond appId, name, clientId and clientSecret
func (c *ApplicationClient) GetByClienId(ctx context.Context, clientId string) (*ApplicationBaseModel, error) {
	filter := fmt.Sprintf("clientId eq %s", clientId)
	excludeSystemPortals := true

	params := GetAllApplicationsParams{
		Filter:               &filter,
		ExcludeSystemPortals: &excludeSystemPortals,
		Attributes:           stringPtr("templateId,clientId,"),
	}

	resp, err := c.apiClient.GetAllApplicationsWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to find application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to find application: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil || resp.JSON200.Applications == nil || len(*resp.JSON200.Applications) == 0 {
		return nil, fmt.Errorf("application with clientId '%s' not found", clientId)
	}

	var targetApp *ApplicationListItem
	for _, app := range *resp.JSON200.Applications {
		if app.ClientId != nil && *app.ClientId == clientId {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		return nil, fmt.Errorf("application with name '%s' not found", clientId)
	}

	return c.getApplicationDetails(ctx, *targetApp.Id)
}

func (c *ApplicationClient) buildSPARequest(name, redirectURL string) (ApplicationModel, error) {
	allowedOrigins, err := extractOrigins(redirectURL)
	if err != nil {
		return ApplicationModel{}, err
	}

	defaultAuthenticationSequenceType := DEFAULT
	defaultClaimDialect := LOCAL

	return ApplicationModel{
		Name: name,
		AdvancedConfigurations: &AdvancedApplicationConfiguration{
			DiscoverableByEndUsers: boolPtr(false),
			SkipLogoutConsent:      boolPtr(true),
			SkipLoginConsent:       boolPtr(true),
		},
		AuthenticationSequence: &AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
			Steps: &[]AuthenticationStepModel{
				{
					Id: 1,
					Options: []Authenticator{
						{
							Idp:           "LOCAL",
							Authenticator: "basic",
						},
					},
				},
			},
		},
		ClaimConfiguration: &ClaimConfiguration{
			Dialect: &defaultClaimDialect,
			RequestedClaims: &[]RequestedClaimConfiguration{
				{
					Claim: Claim{
						Uri: "http://wso2.org/claims/username",
					},
				},
			},
		},
		InboundProtocolConfiguration: &InboundProtocols{
			Oidc: &OpenIDConnectConfiguration{
				AccessToken: &AccessTokenConfiguration{
					ApplicationAccessTokenExpiryInSeconds: int64Ptr(3600),
					BindingType:                           stringPtr("sso-session"),
					RevokeTokensWhenIDPSessionTerminated:  boolPtr(true),
					Type:                                  stringPtr("JWT"),
					UserAccessTokenExpiryInSeconds:        int64Ptr(3600),
					ValidateTokenBinding:                  boolPtr(false),
				},
				GrantTypes:     []string{"authorization_code", "refresh_token"},
				AllowedOrigins: &allowedOrigins,
				CallbackURLs:   &[]string{redirectURL},
				Pkce: &OAuth2PKCEConfiguration{
					Mandatory:                      boolPtr(true),
					SupportPlainTransformAlgorithm: boolPtr(false),
				},
				PublicClient: boolPtr(true),
				RefreshToken: &RefreshTokenConfiguration{
					ExpiryInSeconds:   int64Ptr(86400),
					RenewRefreshToken: boolPtr(true),
				},
			},
		},
		TemplateId: stringPtr("6a90e4b0-fbff-42d7-bfde-1efd98f07cd7"),
		AssociatedRoles: &AssociatedRolesConfig{
			AllowedAudience: APPLICATION,
			Roles:           &[]Role{},
		},
	}, nil
}

func (c *ApplicationClient) buildMobileAppRequest(name, redirectURL string) (ApplicationModel, error) {
	defaultAuthenticationSequenceType := DEFAULT
	return ApplicationModel{
		Name: name,
		AdvancedConfigurations: &AdvancedApplicationConfiguration{
			DiscoverableByEndUsers: boolPtr(false),
			SkipLogoutConsent:      boolPtr(true),
			SkipLoginConsent:       boolPtr(true),
		},
		AuthenticationSequence: &AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
			Steps: &[]AuthenticationStepModel{
				{
					Id: 1,
					Options: []Authenticator{
						{
							Idp:           "LOCAL",
							Authenticator: "basic",
						},
					},
				},
			},
		},
		InboundProtocolConfiguration: &InboundProtocols{
			Oidc: &OpenIDConnectConfiguration{
				AccessToken: &AccessTokenConfiguration{
					ApplicationAccessTokenExpiryInSeconds: int64Ptr(3600),
					BindingType:                           stringPtr("None"),
					RevokeTokensWhenIDPSessionTerminated:  boolPtr(false),
					Type:                                  stringPtr("JWT"),
					UserAccessTokenExpiryInSeconds:        int64Ptr(3600),
					ValidateTokenBinding:                  boolPtr(false),
				},

				GrantTypes:     []string{"authorization_code", "refresh_token"},
				CallbackURLs:   &[]string{redirectURL},
				AllowedOrigins: &[]string{}, // Empty for mobile apps
				Pkce: &OAuth2PKCEConfiguration{
					Mandatory:                      boolPtr(true),
					SupportPlainTransformAlgorithm: boolPtr(false),
				},
				PublicClient: boolPtr(true),
				RefreshToken: &RefreshTokenConfiguration{
					ExpiryInSeconds:   int64Ptr(86400),
					RenewRefreshToken: boolPtr(true),
				},
			},
		},
		TemplateId: stringPtr("mobile-application"),
		AssociatedRoles: &AssociatedRolesConfig{
			AllowedAudience: APPLICATION,
			Roles:           &[]Role{},
		},
	}, nil
}

func (c *ApplicationClient) buildM2MAppRequest(name string) (ApplicationModel, error) {
	defaultAuthenticationSequenceType := DEFAULT
	return ApplicationModel{
		Name: name,
		InboundProtocolConfiguration: &InboundProtocols{
			Oidc: &OpenIDConnectConfiguration{
				GrantTypes:   []string{"client_credentials"},
				PublicClient: boolPtr(false),
			},
		},
		AuthenticationSequence: &AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
		},
		TemplateId: stringPtr("m2m-application"),
		AssociatedRoles: &AssociatedRolesConfig{
			AllowedAudience: APPLICATION,
			Roles:           &[]Role{},
		},
	}, nil
}

func (c *ApplicationClient) processCreatePublicClientAppResponse(ctx context.Context, resp *CreateApplicationResponse, name, redirectURL string) (*ApplicationBaseModel, error) {
	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("failed to create application: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	if resp.HTTPResponse == nil {
		return nil, fmt.Errorf("unexpected empty HTTP response")
	}

	locationHeader := resp.HTTPResponse.Header.Get("Location")
	if locationHeader == "" {
		return nil, fmt.Errorf("location header is missing in the response")
	}

	appID, err := extractApplicationID(locationHeader)
	if err != nil {
		return nil, err
	}

	appDetails, err := c.fetchApplicationDetails(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("created application but failed to fetch details: %w", err)
	}

	return &ApplicationBaseModel{
		Id:               appID,
		Name:             name,
		ClientId:         *appDetails.ClientId,
		RedirectURL:      redirectURL,
		AuthorizedScopes: "openid profile email",
	}, nil
}

func (c *ApplicationClient) processCreateConfidentialClientAppResponse(ctx context.Context, resp *CreateApplicationResponse, name string) (*ApplicationBaseModel, error) {
	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("failed to create application: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	if resp.HTTPResponse == nil {
		return nil, fmt.Errorf("unexpected empty HTTP response")
	}

	locationHeader := resp.HTTPResponse.Header.Get("Location")
	if locationHeader == "" {
		return nil, fmt.Errorf("location header is missing in the response")
	}

	appID, err := extractApplicationID(locationHeader)
	if err != nil {
		return nil, err
	}

	oauthDetails, err := c.fetchInboundOAuthDetails(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OAuth client credentials: %w", err)
	}

	return &ApplicationBaseModel{
		Id:           appID,
		Name:         name,
		ClientId:     *oauthDetails.ClientId,
		ClientSecret: *oauthDetails.ClientSecret,
	}, nil
}

func (c *ApplicationClient) fetchInboundOAuthDetails(ctx context.Context, appID string) (*OpenIDConnectConfiguration, error) {
	resp, err := c.apiClient.GetInboundOAuthConfigurationWithResponse(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth protocol details: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get OAuth protocol details: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected empty response body")
	}

	return resp.JSON200, nil
}

func (c *ApplicationClient) fetchApplicationDetails(ctx context.Context, appID string) (*ApplicationResponseModel, error) {
	resp, err := c.apiClient.GetApplicationWithResponse(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application details: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get application details: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected empty response body")
	}

	return resp.JSON200, nil
}

func (c *ApplicationClient) getApplicationDetails(ctx context.Context, appID string) (*ApplicationBaseModel, error) {
	appDetails, err := c.fetchApplicationDetails(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application details: %w", err)
	}

	result := &ApplicationBaseModel{
		Id:   appID,
		Name: appDetails.Name,
	}

	if appDetails.ClientId != nil {
		result.ClientId = *appDetails.ClientId
	}

	isM2MApp := false
	if appDetails.TemplateId != nil && *appDetails.TemplateId == "m2m-application" {
		isM2MApp = true
	}

	if isM2MApp {
		oauthDetails, err := c.fetchInboundOAuthDetails(ctx, appID)
		if err != nil {
			return nil, fmt.Errorf("failed to get OAuth details: %w", err)
		}

		if oauthDetails.ClientSecret != nil {
			result.ClientSecret = *oauthDetails.ClientSecret
		}
	}

	return result, nil
}

func (c *ApplicationClient) AuthorizeAPI(ctx context.Context, appID string, authorizedAPI AddAuthorizedAPIJSONRequestBody) (*http.Response, error) {
	resp, err := c.apiClient.AddAuthorizedAPIWithResponse(ctx, appID, authorizedAPI)
	if err != nil {
		return nil, fmt.Errorf("Failed to authorize api: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to authorize api: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return &http.Response{}, nil
}
