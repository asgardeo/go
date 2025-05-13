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
	"sort"
	"strings"

	"github.com/asgardeo/go/pkg/application/internal"
	"github.com/asgardeo/go/pkg/authenticator"
	"github.com/asgardeo/go/pkg/claim"
	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/identity_provider"
	"github.com/asgardeo/go/pkg/oidc_scope"
)

const (
	defaultLimit  int = 10
	defaultOffset int = 0
)

// ApplicationClient is a wrapper around the generated client for the Application Management API
type ApplicationClient struct {
	config    *config.ClientConfig
	apiClient *internal.ClientWithResponses
}

// New creates a new Application Management API client
func New(cfg *config.ClientConfig) (*ApplicationClient, error) {
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	apiClient, err := internal.NewClientWithResponses(
		cfg.BaseURL+"/api/server/v1",
		internal.WithHTTPClient(cfg.HTTPClient),
		internal.WithRequestEditorFn(typedAuthEditorFn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create application client: %w", err)
	}

	return &ApplicationClient{
		config:    cfg,
		apiClient: apiClient,
	}, nil
}

// List retrieves a list of applications with pagination support
func (c *ApplicationClient) List(ctx context.Context, limit, offset int) (*ApplicationListResponseModel, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = defaultOffset
	}
	excludeSystemPortals := true
	params := internal.GetAllApplicationsParams{
		Limit:                &limit,
		Offset:               &offset,
		ExcludeSystemPortals: &excludeSystemPortals,
	}
	resp, err := c.apiClient.GetAllApplicationsWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to list applications: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// CreateSinglePageApp creates a new Single Page Application with sensible defaults
func (c *ApplicationClient) CreateSinglePageApp(ctx context.Context, name string, redirectURL string) (*ApplicationBasicInfoResponseModel, error) {
	appRequest, err := c.buildSPARequest(name, redirectURL)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create SPA: %w", err)
	}

	return c.processCreateAppResponse(ctx, resp, name, AppTypeSPA, &redirectURL)
}

// CreateMobileApp creates a new Mobile Application with sensible defaults
func (c *ApplicationClient) CreateMobileApp(ctx context.Context, name string, redirectURL string) (*ApplicationBasicInfoResponseModel, error) {
	appRequest, err := c.buildMobileAppRequest(name, redirectURL)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create mobile app: %w", err)
	}

	return c.processCreateAppResponse(ctx, resp, name, AppTypeMobile, &redirectURL)
}

// CreateM2MApp creates a new Machine-to-Machine (M2M) Application
func (c *ApplicationClient) CreateM2MApp(ctx context.Context, name string) (*ApplicationBasicInfoResponseModel, error) {
	appRequest, err := c.buildM2MAppRequest(name)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create M2M app: %w", err)
	}

	return c.processCreateAppResponse(ctx, resp, name, AppTypeM2M, nil)
}

// CreateWebAppWithSSR creates a new Web Application with Server-Side Rendering support
func (c *ApplicationClient) CreateWebAppWithSSR(ctx context.Context, name string, redirectURL string) (*ApplicationBasicInfoResponseModel, error) {
	appRequest, err := c.buildWebAppWithSSRRequest(name, redirectURL)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.CreateApplicationWithResponse(ctx, nil, appRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSR webapp: %w", err)
	}

	return c.processCreateAppResponse(ctx, resp, name, AppTypeSSRWeb, &redirectURL)
}

// GetByName finds an application by name and returns its details
// todo: improve application details being fetched beyond appId, name, clientId and clientSecret
func (c *ApplicationClient) GetByName(ctx context.Context, name string) (*ApplicationBasicInfoResponseModel, error) {
	filter := fmt.Sprintf("name eq %s", name)
	excludeSystemPortals := true

	params := internal.GetAllApplicationsParams{
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

	var targetApp *internal.ApplicationListItem
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
func (c *ApplicationClient) GetByClienId(ctx context.Context, clientId string) (*ApplicationBasicInfoResponseModel, error) {
	filter := fmt.Sprintf("clientId eq %s", clientId)
	excludeSystemPortals := true

	params := internal.GetAllApplicationsParams{
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

	var targetApp *internal.ApplicationListItem
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

// AuthorizeAPI authorizes an application to access an API with specified scopes
func (c *ApplicationClient) AuthorizeAPI(ctx context.Context, appID string, apiAuthorization AuthorizedAPICreateModel) error {
	resp, err := c.apiClient.AddAuthorizedAPIWithResponse(ctx, appID, apiAuthorization)
	if err != nil {
		return fmt.Errorf("failed to authorize API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to authorize API: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	return nil
}

// GetAuthorizedAPIs retrieves the list of APIs authorized for an application
func (c *ApplicationClient) GetAuthorizedAPIs(ctx context.Context, appID string) (*[]AuthorizedAPIResponseModel, error) {
	resp, err := c.apiClient.GetAuthorizedAPIsWithResponse(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized APIs: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get authorized APIs: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// UpdateBasicInfo updates basic information of an existing application
func (c *ApplicationClient) UpdateBasicInfo(ctx context.Context, appId string, updateModel ApplicationBasicInfoUpdateModel) error {
	patchData := convertBasicInfoUpdateModelToApplicationPatchModel(updateModel)
	resp, err := c.apiClient.PatchApplicationWithResponse(ctx, appId, patchData)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to update application: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	return nil
}

// UpdateOAuthConfig updates allowed OAuth configuration fields for an application
func (c *ApplicationClient) UpdateOAuthConfig(ctx context.Context, applicationId string, config ApplicationOAuthConfigUpdateModel) error {
	resp, err := c.apiClient.GetInboundOAuthConfigurationWithResponse(ctx, applicationId)
	if err != nil {
		return fmt.Errorf("failed to get existing OAuth configuration: %w", err)
	}

	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return fmt.Errorf("failed to get existing OAuth configuration: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}

	updatedConfig := *resp.JSON200

	if updatedConfig.AccessToken == nil {
		updatedConfig.AccessToken = &internal.AccessTokenConfiguration{}
	}

	if config.AccessTokenAttributes != nil {
		updatedConfig.AccessToken.AccessTokenAttributes = config.AccessTokenAttributes
	}

	// todo: Allow application access token expiry only for M2M and SSR web apps
	if config.ApplicationAccessTokenExpiryInSeconds != nil {
		updatedConfig.AccessToken.ApplicationAccessTokenExpiryInSeconds = config.ApplicationAccessTokenExpiryInSeconds
	}

	// todo: Allow UserAccessTokenExpiryInSeconds only for mobile, SPA, web apps
	if config.UserAccessTokenExpiryInSeconds != nil {
		updatedConfig.AccessToken.UserAccessTokenExpiryInSeconds = config.UserAccessTokenExpiryInSeconds
	}

	// todo: Allow CORS origins only or mobile, SPA, web app
	if config.AllowedOrigins != nil {
		updatedConfig.AllowedOrigins = config.AllowedOrigins
	}

	if config.CallbackURLs != nil {
		if len(*config.CallbackURLs) == 1 {
			// If there's only one callback URL, use it directly without "regexp="
			updatedConfig.CallbackURLs = config.CallbackURLs
		} else {
			// Construct "regexp=(callback1|callback2|...|callbackN)" for multiple callback URLs
			var callbackRegex string
			for i, callback := range *config.CallbackURLs {
				if i == 0 {
					callbackRegex = fmt.Sprintf("regexp=(%s", callback)
				} else {
					callbackRegex = fmt.Sprintf("%s|%s", callbackRegex, callback)
				}
			}
			if len(*config.CallbackURLs) > 0 {
				callbackRegex = fmt.Sprintf("%s)", callbackRegex)
			} else {
				callbackRegex = ""
			}

			updatedConfig.CallbackURLs = &[]string{callbackRegex}
		}
	}

	if config.Logout != nil {
		updatedConfig.Logout = config.Logout
	}

	if config.RefreshTokenExpiryInSeconds != nil {
		if updatedConfig.RefreshToken == nil {
			updatedConfig.RefreshToken = &internal.RefreshTokenConfiguration{}
		}
		updatedConfig.RefreshToken.ExpiryInSeconds = config.RefreshTokenExpiryInSeconds
	}

	// Update the configuration
	updateResp, err := c.apiClient.UpdateInboundOAuthConfigurationWithResponse(ctx, applicationId, updatedConfig)
	if err != nil {
		return fmt.Errorf("failed to update OAuth configuration: %w", err)
	}

	if updateResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to update OAuth configuration: status %d, body: %s",
			updateResp.StatusCode(), string(updateResp.Body))
	}

	return nil
}

// UpdateClaimConfig updates the claim configuration of an existing application
func (c *ApplicationClient) UpdateClaimConfig(ctx context.Context, appId string, claimConfigUpdateModel ApplicationClaimConfigurationUpdateModel) error {
	patchData := convertClaimConfigUpdateModelToApplicationPatchModel(claimConfigUpdateModel)
	resp, err := c.apiClient.PatchApplicationWithResponse(ctx, appId, patchData)
	if err != nil {
		return fmt.Errorf("failed to update claim configuration: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to update claim configuration: status %d, body: %s",
			resp.StatusCode(), string(resp.Body))
	}
	return nil
}

// UpdateLoginFlow updates the login flow of an existing application.
func (c *ApplicationClient) UpdateLoginFlow(ctx context.Context, appId string, loginFlowUpdateRequest LoginFlowUpdateModel) error {
	authenticationSequence := internal.ApplicationPatchModel{
		AuthenticationSequence: &loginFlowUpdateRequest,
	}
	resp, err := c.apiClient.PatchApplicationWithResponse(ctx, appId, authenticationSequence)
	if err != nil {
		return fmt.Errorf("failed to update login flow: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to update login flow: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return nil
}

// GenerateLoginFlow initiates the login flow generation process for an application.
func (c *ApplicationClient) GenerateLoginFlow(ctx context.Context, prompt string) (*LoginFlowGenerateResponseModel, error) {

	availableAuthenticators, err := c.buildAvailableAuthenticators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build available authenticators: %w", err)
	}

	userClaims, err := c.buildUserClaimList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build user claims: %w", err)
	}
	loginFlowGenerateRequest := internal.LoginFlowGenerateRequest{
		AvailableAuthenticators: &availableAuthenticators,
		UserClaims:              &userClaims,
		UserQuery:               &prompt,
	}
	resp, err := c.apiClient.GenerateLoginFlowWithResponse(ctx, loginFlowGenerateRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate login flow: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to generate login flow: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// GetLoginFlowGenerationStatus retrieves the status of the login flow generation process.
func (c *ApplicationClient) GetLoginFlowGenerationStatus(ctx context.Context, flowId string) (*LoginFlowStatusResponseModel, error) {
	resp, err := c.apiClient.GetLoginFlowGenerationStatusWithResponse(ctx, flowId)
	if err != nil {
		return nil, fmt.Errorf("failed to get login flow generation status: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get login flow generation status: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

// GetLoginFlowGenerationResult retrieves the result of the login flow generation process.
func (c *ApplicationClient) GetLoginFlowGenerationResult(ctx context.Context, flowId string) (*LoginFlowResultResponseModel, error) {
	resp, err := c.apiClient.GetLoginFlowGenerationResultWithResponse(ctx, flowId)
	if err != nil {
		return nil, fmt.Errorf("failed to get login flow generation result: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get login flow generation result: status %d, body: %s", resp.StatusCode(), string(resp.Body))
	}
	loginFlowResultResponse := convertToLoginFlowResultResponseModel(*resp.JSON200)
	return &loginFlowResultResponse, nil
}

func (c *ApplicationClient) buildSPARequest(name, redirectURL string) (internal.ApplicationModel, error) {
	allowedOrigins, err := extractOrigins(redirectURL)
	if err != nil {
		return internal.ApplicationModel{}, err
	}

	defaultAuthenticationSequenceType := internal.DEFAULT
	defaultClaimDialect := internal.LOCAL

	return internal.ApplicationModel{
		Name: name,
		AdvancedConfigurations: &internal.AdvancedApplicationConfiguration{
			DiscoverableByEndUsers: boolPtr(false),
			SkipLogoutConsent:      boolPtr(true),
			SkipLoginConsent:       boolPtr(true),
		},
		AuthenticationSequence: &internal.AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
			Steps: &[]internal.AuthenticationStepModel{
				{
					Id: 1,
					Options: []internal.Authenticator{
						{
							Idp:           "LOCAL",
							Authenticator: "basic",
						},
					},
				},
			},
		},
		ClaimConfiguration: &internal.ClaimConfiguration{
			Dialect: &defaultClaimDialect,
			RequestedClaims: &[]internal.RequestedClaimConfiguration{
				{
					Claim: internal.Claim{
						Uri: "http://wso2.org/claims/username",
					},
				},
			},
		},
		InboundProtocolConfiguration: &internal.InboundProtocols{
			Oidc: &internal.OpenIDConnectConfiguration{
				AccessToken: &internal.AccessTokenConfiguration{
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
				Pkce: &internal.OAuth2PKCEConfiguration{
					Mandatory:                      boolPtr(true),
					SupportPlainTransformAlgorithm: boolPtr(false),
				},
				PublicClient: boolPtr(true),
				RefreshToken: &internal.RefreshTokenConfiguration{
					ExpiryInSeconds:   int64Ptr(86400),
					RenewRefreshToken: boolPtr(true),
				},
			},
		},
		TemplateId: stringPtr("6a90e4b0-fbff-42d7-bfde-1efd98f07cd7"),
		AssociatedRoles: &internal.AssociatedRolesConfig{
			AllowedAudience: internal.APPLICATION,
			Roles:           &[]internal.Role{},
		},
	}, nil
}

func (c *ApplicationClient) buildMobileAppRequest(name, redirectURL string) (internal.ApplicationModel, error) {
	defaultAuthenticationSequenceType := internal.DEFAULT
	return internal.ApplicationModel{
		Name: name,
		AdvancedConfigurations: &internal.AdvancedApplicationConfiguration{
			DiscoverableByEndUsers: boolPtr(false),
			SkipLogoutConsent:      boolPtr(true),
			SkipLoginConsent:       boolPtr(true),
		},
		AuthenticationSequence: &internal.AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
			Steps: &[]internal.AuthenticationStepModel{
				{
					Id: 1,
					Options: []internal.Authenticator{
						{
							Idp:           "LOCAL",
							Authenticator: "basic",
						},
					},
				},
			},
		},
		InboundProtocolConfiguration: &internal.InboundProtocols{
			Oidc: &internal.OpenIDConnectConfiguration{
				AccessToken: &internal.AccessTokenConfiguration{
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
				Pkce: &internal.OAuth2PKCEConfiguration{
					Mandatory:                      boolPtr(true),
					SupportPlainTransformAlgorithm: boolPtr(false),
				},
				PublicClient: boolPtr(true),
				RefreshToken: &internal.RefreshTokenConfiguration{
					ExpiryInSeconds:   int64Ptr(86400),
					RenewRefreshToken: boolPtr(true),
				},
			},
		},
		TemplateId: stringPtr("mobile-application"),
		AssociatedRoles: &internal.AssociatedRolesConfig{
			AllowedAudience: internal.APPLICATION,
			Roles:           &[]internal.Role{},
		},
	}, nil
}

func (c *ApplicationClient) buildM2MAppRequest(name string) (internal.ApplicationModel, error) {
	defaultAuthenticationSequenceType := internal.DEFAULT
	return internal.ApplicationModel{
		Name: name,
		InboundProtocolConfiguration: &internal.InboundProtocols{
			Oidc: &internal.OpenIDConnectConfiguration{
				GrantTypes:   []string{"client_credentials"},
				PublicClient: boolPtr(false),
			},
		},
		AuthenticationSequence: &internal.AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
		},
		TemplateId: stringPtr("m2m-application"),
		AssociatedRoles: &internal.AssociatedRolesConfig{
			AllowedAudience: internal.APPLICATION,
			Roles:           &[]internal.Role{},
		},
	}, nil
}

func (c *ApplicationClient) buildWebAppWithSSRRequest(name, redirectURL string) (internal.ApplicationModel, error) {
	defaultAuthenticationSequenceType := internal.DEFAULT
	defaultClaimDialect := internal.LOCAL

	return internal.ApplicationModel{
		Name: name,
		AdvancedConfigurations: &internal.AdvancedApplicationConfiguration{
			DiscoverableByEndUsers: boolPtr(false),
			SkipLogoutConsent:      boolPtr(true),
			SkipLoginConsent:       boolPtr(true),
		},
		AuthenticationSequence: &internal.AuthenticationSequence{
			Type: &defaultAuthenticationSequenceType,
			Steps: &[]internal.AuthenticationStepModel{
				{
					Id: 1,
					Options: []internal.Authenticator{
						{
							Idp:           "LOCAL",
							Authenticator: "basic",
						},
					},
				},
			},
		},
		ClaimConfiguration: &internal.ClaimConfiguration{
			Dialect: &defaultClaimDialect,
			RequestedClaims: &[]internal.RequestedClaimConfiguration{
				{
					Claim: internal.Claim{
						Uri: "http://wso2.org/claims/username",
					},
				},
			},
		},
		InboundProtocolConfiguration: &internal.InboundProtocols{
			Oidc: &internal.OpenIDConnectConfiguration{
				GrantTypes:     []string{"authorization_code"},
				CallbackURLs:   &[]string{redirectURL},
				AllowedOrigins: &[]string{},
				PublicClient:   boolPtr(false),
				RefreshToken: &internal.RefreshTokenConfiguration{
					ExpiryInSeconds: int64Ptr(86400),
				},
				// Including access token configuration appropriate for SSR webapps
				AccessToken: &internal.AccessTokenConfiguration{
					ApplicationAccessTokenExpiryInSeconds: int64Ptr(3600),
					BindingType:                           stringPtr("cookie"),
					RevokeTokensWhenIDPSessionTerminated:  boolPtr(true),
					Type:                                  stringPtr("JWT"),
					UserAccessTokenExpiryInSeconds:        int64Ptr(3600),
				},
			},
		},
		TemplateId: stringPtr("b9c5e11e-fc78-484b-9bec-015d247561b8"), // Web application template
		AssociatedRoles: &internal.AssociatedRolesConfig{
			AllowedAudience: internal.APPLICATION,
			Roles:           &[]internal.Role{},
		},
	}, nil
}

func (c *ApplicationClient) processCreateAppResponse(ctx context.Context, resp *internal.CreateApplicationResponse, name string, appType AppType, redirectURL *string) (*ApplicationBasicInfoResponseModel, error) {
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

	if appType == AppTypeSPA || appType == AppTypeMobile {
		appDetails, err := c.fetchApplicationDetails(ctx, appID)
		if err != nil {
			return nil, fmt.Errorf("created application but failed to fetch details: %w", err)
		}

		return &ApplicationBasicInfoResponseModel{
			Id:               appID,
			Name:             name,
			ClientId:         *appDetails.ClientId,
			RedirectURL:      *redirectURL,
			AuthorizedScopes: "openid profile email",
			AppType:          appType,
		}, nil
	}

	oauthDetails, err := c.fetchInboundOAuthDetails(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OAuth client credentials: %w", err)
	}

	if appType == AppTypeSSRWeb {
		return &ApplicationBasicInfoResponseModel{
			Id:               appID,
			Name:             name,
			ClientId:         *oauthDetails.ClientId,
			ClientSecret:     *oauthDetails.ClientSecret,
			AuthorizedScopes: "openid profile email",
			AppType:          appType,
		}, nil
	}

	if appType == AppTypeM2M {

		return &ApplicationBasicInfoResponseModel{
			Id:               appID,
			Name:             name,
			ClientId:         *oauthDetails.ClientId,
			ClientSecret:     *oauthDetails.ClientSecret,
			AuthorizedScopes: "openid profile email",
			AppType:          appType,
		}, nil
	}

	return nil, err
}

func (c *ApplicationClient) fetchInboundOAuthDetails(ctx context.Context, appID string) (*internal.OpenIDConnectConfiguration, error) {
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

func (c *ApplicationClient) fetchApplicationDetails(ctx context.Context, appID string) (*internal.ApplicationResponseModel, error) {
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

func (c *ApplicationClient) getApplicationDetails(ctx context.Context, appID string) (*ApplicationBasicInfoResponseModel, error) {
	appDetails, err := c.fetchApplicationDetails(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application details: %w", err)
	}

	oauthDetails, err := c.fetchInboundOAuthDetails(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth details: %w", err)
	}

	authorizedAPIs, err := c.GetAuthorizedAPIs(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized APIs: %w", err)
	}

	scopeSet := make(map[string]struct{})
	if authorizedAPIs != nil {
		for _, api := range *authorizedAPIs {
			if api.AuthorizedScopes != nil {
				for _, scope := range *api.AuthorizedScopes {
					scopeSet[*scope.Name] = struct{}{}
				}
			}
		}
	}

	var scopes []string
	for scope := range scopeSet {
		scopes = append(scopes, scope)
	}

	result := &ApplicationBasicInfoResponseModel{
		Id:   appID,
		Name: appDetails.Name,
	}

	if appDetails.ClientId != nil {
		result.ClientId = *appDetails.ClientId
	}

	appType, err := determineAppType(appDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to determine application type: %w", err)
	}

	if appType == AppTypeM2M || appType == AppTypeSSRWeb {
		if oauthDetails.ClientSecret != nil {
			result.ClientSecret = *oauthDetails.ClientSecret
		}
	}

	if appType == AppTypeM2M {
		if len(scopes) > 0 {
			result.AuthorizedScopes = strings.Join(scopes, " ")
		}
	}

	if appType == AppTypeSPA || appType == AppTypeMobile || appType == AppTypeSSRWeb {
		if oauthDetails.CallbackURLs != nil && len(*oauthDetails.CallbackURLs) > 0 {
			firstCallback := (*oauthDetails.CallbackURLs)[0]
			if strings.HasPrefix(firstCallback, "regexp=") {
				// Extract URLs from the "regexp=(url1|url2|...|urlN)" format
				regexContent := strings.TrimPrefix(firstCallback, "regexp=")
				urls := strings.Split(strings.Trim(regexContent, "()"), "|")
				result.RedirectURL = strings.Join(urls, ",")
			} else {
				result.RedirectURL = firstCallback
			}
		}

		authorizedOIDCScopeList, err := c.getAuthorizedOIDCScopes(ctx, appDetails.ClaimConfiguration)
		if err != nil {
			return nil, fmt.Errorf("failed to get authorized OIDC scopes: %w", err)
		}
		authorizedScopes := append(authorizedOIDCScopeList, scopes...)
		result.AuthorizedScopes = strings.Join(authorizedScopes, " ")
	}

	return result, nil
}

func determineAppType(appDetails *internal.ApplicationResponseModel) (AppType, error) {
	if appDetails.TemplateId != nil {
		switch *appDetails.TemplateId {
		case "6a90e4b0-fbff-42d7-bfde-1efd98f07cd7": // SPA Template ID
			return AppTypeSPA, nil
		case "mobile-application": // Mobile App Template ID
			return AppTypeMobile, nil
		case "m2m-application": // M2M Template ID
			return AppTypeM2M, nil
		case "b9c5e11e-fc78-484b-9bec-015d247561b8": // Web App with SSR Template ID
			return AppTypeSSRWeb, nil
		}
	}
	return "", fmt.Errorf("unknown application type")
}

func (c *ApplicationClient) buildAvailableAuthenticators(ctx context.Context) (map[string]interface{}, error) {
	authenticatorClient, err := authenticator.New(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator client: %w", err)
	}
	localAuthenticators, err := authenticatorClient.ListLocalAuthenticators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list local authenticators: %w", err)
	}

	var moderatedAuthenticators []interface{}
	var secondFactorAuthenticators []interface{}
	var recoveryAuthenticators []interface{}
	for _, localAuthenticator := range *localAuthenticators {
		description := ""
		if localAuthenticator.Description != nil {
			description = *localAuthenticator.Description
		}
		authenticatorData := map[string]interface{}{
			"description": description,
			"idp":         *localAuthenticator.Type,
			"name":        *localAuthenticator.Name,
		}

		if internal.LocalAuthenticatorIDs.BackupCode == *localAuthenticator.Id {
			recoveryAuthenticators = append(recoveryAuthenticators, authenticatorData)
		} else if _, exists := internal.SecondFactorAuthenticatorIDs[*localAuthenticator.Id]; exists {
			secondFactorAuthenticators = append(secondFactorAuthenticators, authenticatorData)
		} else {
			moderatedAuthenticators = append(moderatedAuthenticators, authenticatorData)
		}
	}

	identityProviderClient, err := identity_provider.New(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity provider client: %w", err)
	}
	requiredAttributes := "federatedAuthenticators"
	getIDPListParams := identity_provider.IdentityProviderListParamsModel{
		RequiredAttributes: &requiredAttributes,
	}
	idpList, err := identityProviderClient.List(ctx, &getIDPListParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list identity providers: %w", err)
	}

	var enterpriseAuthenticators []interface{}
	var socialAuthenticators []interface{}
	for _, idp := range *idpList.IdentityProviders {
		if idp.Name == nil {
			// Skip the IdP if the name is not available
			continue
		}

		federatedAuthenticators := idp.FederatedAuthenticators
		if federatedAuthenticators == nil {
			// Skip the IdP if federated authenticator list is not available
			continue
		}

		defaultAuthenticatorId := federatedAuthenticators.DefaultAuthenticatorId
		if defaultAuthenticatorId == nil {
			// Skip the IdP if default authenticator ID is not available
			continue
		}

		description := ""
		if idp.Description != nil {
			description = *idp.Description
		}

		idpAuthenticatorName := ""
		for _, authenticator := range *federatedAuthenticators.Authenticators {
			if authenticator.AuthenticatorId != nil && *authenticator.AuthenticatorId == *defaultAuthenticatorId {
				if authenticator.Name != nil {
					idpAuthenticatorName = *authenticator.Name
				}
				break
			}
		}

		if idpAuthenticatorName == "" {
			// Skip the IdP if the authenticator name is not found
			continue
		}

		authenticatorData := map[string]interface{}{
			"description": description,
			"idp":         *idp.Name,
			"name":        idpAuthenticatorName,
		}

		if _, exists := internal.SocialAuthenticatorIDs[*defaultAuthenticatorId]; exists {
			socialAuthenticators = append(socialAuthenticators, authenticatorData)
		} else {
			enterpriseAuthenticators = append(enterpriseAuthenticators, authenticatorData)
		}
	}

	availableAuthenticators := map[string]interface{}{
		"enterprise":   enterpriseAuthenticators,
		"local":        moderatedAuthenticators,
		"recovery":     recoveryAuthenticators,
		"secondFactor": secondFactorAuthenticators,
		"social":       socialAuthenticators,
	}
	return availableAuthenticators, nil
}

func (c *ApplicationClient) buildUserClaimList(ctx context.Context) ([]map[string]interface{}, error) {
	claimClient, err := claim.New(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create claim client: %w", err)
	}
	excludeHiddenClaims := true
	listLocalClaimsParams := claim.LocalClaimListParamsModel{
		ExcludeHiddenClaims: &excludeHiddenClaims,
	}
	claims, err := claimClient.ListLocalClaims(ctx, &listLocalClaimsParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list local claims: %w", err)
	}

	if claims == nil || *claims == nil {
		return nil, fmt.Errorf("failed to list local claims: empty response")
	}

	userClaims := []map[string]interface{}{}
	for _, claim := range *claims {
		if claim.ClaimURI != nil && claim.Description != nil {
			userClaims = append(userClaims, map[string]interface{}{
				"claimURI":    *claim.ClaimURI,
				"description": *claim.Description,
			})
		}
	}

	return userClaims, nil
}

func (c *ApplicationClient) getAuthorizedOIDCScopes(ctx context.Context, claimConfig *internal.ClaimConfiguration) ([]string, error) {
	if claimConfig == nil || *claimConfig.RequestedClaims == nil {
		// If no requested claims are found, return only the openid scope
		return []string{"openid"}, nil
	}
	requestedClaims := *claimConfig.RequestedClaims

	oidcScopeClient, err := oidc_scope.New(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC scope client: %w", err)
	}
	oidcScopeListResponse, err := oidcScopeClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list OIDC scopes: %w", err)
	}
	if oidcScopeListResponse == nil || *oidcScopeListResponse == nil {
		// If no OIDC scopes are found, return only the openid scope
		return []string{"openid"}, nil
	}
	oidcScopeList := *oidcScopeListResponse

	claimClient, err := claim.New(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create claim client: %w", err)
	}
	oidcClaimListResponse, err := claimClient.ListOIDCClaims(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list OIDC claims: %w", err)
	}
	if oidcClaimListResponse == nil || *oidcClaimListResponse == nil {
		// If no OIDC claims are found, return only the openid scope
		return []string{"openid"}, nil
	}
	oidcClaimList := *oidcClaimListResponse

	authorizedOIDCScopes := make(map[string]bool)
	localToOIDCClaimMap := make(map[string]string)
	for _, oidcClaim := range oidcClaimList {
		if oidcClaim.MappedLocalClaimURI != nil && oidcClaim.ClaimURI != nil {
			localToOIDCClaimMap[*oidcClaim.MappedLocalClaimURI] = *oidcClaim.ClaimURI
		}
	}

	for _, requestedClaim := range requestedClaims {
		if oidcClaimURI, ok := localToOIDCClaimMap[requestedClaim.Claim.Uri]; ok {
			for _, oidcScope := range oidcScopeList {
				for _, claimInScope := range oidcScope.Claims {
					if claimInScope == oidcClaimURI && oidcScope.Name != "openid" {
						authorizedOIDCScopes[oidcScope.Name] = true
						break
					}
				}
			}
		}
	}

	result := make([]string, 0, len(authorizedOIDCScopes)+1)
	// Add "openid" scope to the result as the first element
	result = append(result, "openid")
	for scope := range authorizedOIDCScopes {
		result = append(result, scope)
	}
	if len(result) > 1 {
		sort.Strings(result[1:])
	}

	return result, nil
}
