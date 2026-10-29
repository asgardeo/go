/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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

package self_service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/self_service/internal/application"
	"github.com/asgardeo/go/pkg/self_service/internal/association"
	"github.com/asgardeo/go/pkg/self_service/internal/authorized_apps"
	"github.com/asgardeo/go/pkg/self_service/internal/backup_code"
	"github.com/asgardeo/go/pkg/self_service/internal/consent"
	"github.com/asgardeo/go/pkg/self_service/internal/fido"
	"github.com/asgardeo/go/pkg/self_service/internal/mfa"
	"github.com/asgardeo/go/pkg/self_service/internal/push"
	"github.com/asgardeo/go/pkg/self_service/internal/session"
	"github.com/asgardeo/go/pkg/self_service/internal/totp"
)

// SelfServiceClient manages the account of the user the access token belongs to
type SelfServiceClient struct {
	config         *config.ClientConfig
	sessions       *session.ClientWithResponses
	associations   *association.ClientWithResponses
	authorizedApps *authorized_apps.ClientWithResponses
	backupCodes    *backup_code.ClientWithResponses
	mfa            *mfa.ClientWithResponses
	totp           *totp.ClientWithResponses
	push           *push.ClientWithResponses
	fido           *fido.ClientWithResponses
	applications   *application.ClientWithResponses
	consents       *consent.ClientWithResponses
}

// New creates a new self-service client
func New(cfg *config.ClientConfig) (*SelfServiceClient, error) {
	authEditorFn := common.CreateAuthRequestEditorFunc(cfg)

	typedAuthEditorFn := func(ctx context.Context, req *http.Request) error {
		editorFn := authEditorFn.(func(context.Context, *http.Request) error)
		return editorFn(ctx, req)
	}

	usersV1 := cfg.BaseURL + "/api/users/v1"
	usersV2 := cfg.BaseURL + "/api/users/v2"
	c := &SelfServiceClient{config: cfg}
	var err error
	if c.sessions, err = session.NewClientWithResponses(usersV1,
		session.WithHTTPClient(cfg.HTTPClient), session.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create session client: %w", err)
	}
	if c.associations, err = association.NewClientWithResponses(usersV1,
		association.WithHTTPClient(cfg.HTTPClient), association.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create association client: %w", err)
	}
	if c.authorizedApps, err = authorized_apps.NewClientWithResponses(usersV2,
		authorized_apps.WithHTTPClient(cfg.HTTPClient), authorized_apps.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create authorized apps client: %w", err)
	}
	if c.backupCodes, err = backup_code.NewClientWithResponses(usersV1,
		backup_code.WithHTTPClient(cfg.HTTPClient), backup_code.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create backup code client: %w", err)
	}
	if c.mfa, err = mfa.NewClientWithResponses(usersV1,
		mfa.WithHTTPClient(cfg.HTTPClient), mfa.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create MFA client: %w", err)
	}
	if c.totp, err = totp.NewClientWithResponses(usersV1,
		totp.WithHTTPClient(cfg.HTTPClient), totp.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create TOTP client: %w", err)
	}
	if c.push, err = push.NewClientWithResponses(usersV1+"/me/push",
		push.WithHTTPClient(cfg.HTTPClient), push.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create push device client: %w", err)
	}
	if c.fido, err = fido.NewClientWithResponses(usersV2+"/me/webauthn",
		fido.WithHTTPClient(cfg.HTTPClient), fido.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create FIDO client: %w", err)
	}
	if c.applications, err = application.NewClientWithResponses(usersV1,
		application.WithHTTPClient(cfg.HTTPClient), application.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create application client: %w", err)
	}
	if c.consents, err = consent.NewClientWithResponses(usersV1,
		consent.WithHTTPClient(cfg.HTTPClient), consent.WithRequestEditorFn(typedAuthEditorFn)); err != nil {
		return nil, fmt.Errorf("failed to create consent client: %w", err)
	}
	return c, nil
}

// ListSessions retrieves the user's active sessions
func (c *SelfServiceClient) ListSessions(ctx context.Context) (*SessionListResponseModel, error) {
	resp, err := c.sessions.GetSessionsOfLoggedInUserWithResponse(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	if err := common.CheckStatus("list sessions", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// TerminateSession terminates a session of the user
func (c *SelfServiceClient) TerminateSession(ctx context.Context, sessionID string) error {
	resp, err := c.sessions.TerminateSessionByLoggedInUserWithResponse(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to terminate session: %w", err)
	}
	return common.CheckStatus("terminate session", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// TerminateAllSessions terminates all sessions of the user
func (c *SelfServiceClient) TerminateAllSessions(ctx context.Context) error {
	resp, err := c.sessions.TerminateSessionsByLoggedInUserWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to terminate sessions: %w", err)
	}
	return common.CheckStatus("terminate sessions", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// GetEnabledMFAAuthenticators retrieves the user's enabled MFA authenticators
func (c *SelfServiceClient) GetEnabledMFAAuthenticators(ctx context.Context) (*EnabledAuthenticatorsResponseModel, error) {
	resp, err := c.mfa.GetMeMfaAuthenticatorsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled MFA authenticators: %w", err)
	}
	if err := common.CheckStatus("get enabled MFA authenticators", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// GetRemainingBackupCodes retrieves the count of the user's remaining backup codes
func (c *SelfServiceClient) GetRemainingBackupCodes(ctx context.Context) (*RemainingBackupCodesResponseModel, error) {
	resp, err := c.backupCodes.GetMeBackupCodesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get remaining backup codes: %w", err)
	}
	if err := common.CheckStatus("get remaining backup codes", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// DeleteBackupCodes deletes the user's backup codes
func (c *SelfServiceClient) DeleteBackupCodes(ctx context.Context) error {
	resp, err := c.backupCodes.DeleteMeBackupCodesWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete backup codes: %w", err)
	}
	return common.CheckStatus("delete backup codes", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// DeleteTOTP removes the user's TOTP authenticator
func (c *SelfServiceClient) DeleteTOTP(ctx context.Context) error {
	resp, err := c.totp.DeleteMeTotpWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete TOTP: %w", err)
	}
	return common.CheckStatus("delete TOTP", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// ListPushDevices retrieves the user's push devices
func (c *SelfServiceClient) ListPushDevices(ctx context.Context) ([]PushDeviceModel, error) {
	resp, err := c.push.GetDevicesWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list push devices: %w", err)
	}
	if err := common.CheckStatus("list push devices", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []PushDeviceModel{}, nil
	}
	return *resp.JSON200, nil
}

// DeletePushDevice removes a push device
func (c *SelfServiceClient) DeletePushDevice(ctx context.Context, deviceID string) error {
	resp, err := c.push.DeleteDeviceByIdWithResponse(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to delete push device: %w", err)
	}
	return common.CheckStatus("delete push device", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// ListPasskeys retrieves the user's passkeys
func (c *SelfServiceClient) ListPasskeys(ctx context.Context) ([]PasskeyModel, error) {
	resp, err := c.fido.GetWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list passkeys: %w", err)
	}
	if err := common.CheckStatus("list passkeys", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []PasskeyModel{}, nil
	}
	return *resp.JSON200, nil
}

// RenamePasskey changes the display name of a passkey
func (c *SelfServiceClient) RenamePasskey(ctx context.Context, credentialID, displayName string) error {
	body := fido.PatchRequest{
		{Operation: fido.REPLACE, Path: "/displayName", Value: &displayName},
	}
	resp, err := c.fido.PatchCredentialIdWithResponse(ctx, credentialID, body)
	if err != nil {
		return fmt.Errorf("failed to rename passkey: %w", err)
	}
	return common.CheckStatus("rename passkey", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// DeletePasskey removes a passkey
func (c *SelfServiceClient) DeletePasskey(ctx context.Context, credentialID string) error {
	resp, err := c.fido.DeleteCredentialIdWithResponse(ctx, credentialID)
	if err != nil {
		return fmt.Errorf("failed to delete passkey: %w", err)
	}
	return common.CheckStatus("delete passkey", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// ListAssociations retrieves the user's linked local accounts
func (c *SelfServiceClient) ListAssociations(ctx context.Context) ([]AssociatedUserModel, error) {
	resp, err := c.associations.GetMeAssociationsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list associations: %w", err)
	}
	if err := common.CheckStatus("list associations", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []AssociatedUserModel{}, nil
	}
	return *resp.JSON200, nil
}

// DeleteAssociation unlinks a local account
func (c *SelfServiceClient) DeleteAssociation(ctx context.Context, associatedUserID string) error {
	resp, err := c.associations.DeleteMeAssociationsAssociatedUserIdWithResponse(ctx, associatedUserID)
	if err != nil {
		return fmt.Errorf("failed to delete association: %w", err)
	}
	return common.CheckStatus("delete association", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// DeleteAllAssociations unlinks all local accounts
func (c *SelfServiceClient) DeleteAllAssociations(ctx context.Context) error {
	resp, err := c.associations.DeleteMeAssociationsWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete associations: %w", err)
	}
	return common.CheckStatus("delete associations", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// ListFederatedAssociations retrieves the user's linked federated accounts
func (c *SelfServiceClient) ListFederatedAssociations(ctx context.Context) ([]FederatedAssociationModel, error) {
	resp, err := c.associations.GetMeFederatedAssociationsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list federated associations: %w", err)
	}
	if err := common.CheckStatus("list federated associations", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []FederatedAssociationModel{}, nil
	}
	return *resp.JSON200, nil
}

// DeleteFederatedAssociation unlinks a federated account
func (c *SelfServiceClient) DeleteFederatedAssociation(ctx context.Context, id string) error {
	resp, err := c.associations.DeleteMeFederatedAssociationsIdWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete federated association: %w", err)
	}
	return common.CheckStatus("delete federated association", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// DeleteAllFederatedAssociations unlinks all federated accounts
func (c *SelfServiceClient) DeleteAllFederatedAssociations(ctx context.Context) error {
	resp, err := c.associations.DeleteMeFederatedAssociationsWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete federated associations: %w", err)
	}
	return common.CheckStatus("delete federated associations", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// ListApplications retrieves the applications the user can access
func (c *SelfServiceClient) ListApplications(ctx context.Context, params *ApplicationListParamsModel) (*ApplicationListResponseModel, error) {
	resp, err := c.applications.GetApplicationsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	if err := common.CheckStatus("list applications", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// GetApplication retrieves an application the user can access
func (c *SelfServiceClient) GetApplication(ctx context.Context, applicationID string) (*ApplicationModel, error) {
	resp, err := c.applications.GetApplicationWithResponse(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	if err := common.CheckStatus("get application", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// ListAuthorizedApps retrieves the user's authorized OAuth applications
func (c *SelfServiceClient) ListAuthorizedApps(ctx context.Context) ([]AuthorizedAppModel, error) {
	resp, err := c.authorizedApps.ListLoggedInUserAuthorizedAppsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list authorized apps: %w", err)
	}
	if err := common.CheckStatus("list authorized apps", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []AuthorizedAppModel{}, nil
	}
	return *resp.JSON200, nil
}

// GetAuthorizedApp retrieves an authorized OAuth application
func (c *SelfServiceClient) GetAuthorizedApp(ctx context.Context, applicationID string) (*AuthorizedAppModel, error) {
	resp, err := c.authorizedApps.GetLoggedInUserAuthorizedAppByAppIdWithResponse(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorized app: %w", err)
	}
	if err := common.CheckStatus("get authorized app", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// RevokeAuthorizedApp revokes an authorized OAuth application
func (c *SelfServiceClient) RevokeAuthorizedApp(ctx context.Context, applicationID string) error {
	resp, err := c.authorizedApps.DeleteLoggedInUserAuthorizedAppByAppIdWithResponse(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to revoke authorized app: %w", err)
	}
	return common.CheckStatus("revoke authorized app", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}

// ListConsents retrieves the user's consents
func (c *SelfServiceClient) ListConsents(ctx context.Context, params *ConsentListParamsModel) ([]ConsentSummaryModel, error) {
	resp, err := c.consents.ListConsentsOfLoggedInUserWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list consents: %w", err)
	}
	if err := common.CheckStatus("list consents", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return []ConsentSummaryModel{}, nil
	}
	return *resp.JSON200, nil
}

// GetConsent retrieves a consent of the user
func (c *SelfServiceClient) GetConsent(ctx context.Context, consentID string) (*ConsentModel, error) {
	resp, err := c.consents.GetConsentOfLoggedInUserWithResponse(ctx, consentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consent: %w", err)
	}
	if err := common.CheckStatus("get consent", resp.StatusCode(), resp.Body, http.StatusOK); err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

// RevokeConsent revokes a consent of the user
func (c *SelfServiceClient) RevokeConsent(ctx context.Context, consentID string) error {
	resp, err := c.consents.RevokeConsentOfLoggedInUserWithResponse(ctx, consentID)
	if err != nil {
		return fmt.Errorf("failed to revoke consent: %w", err)
	}
	return common.CheckStatus("revoke consent", resp.StatusCode(), resp.Body, http.StatusOK, http.StatusNoContent)
}
