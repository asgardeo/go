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
	"github.com/asgardeo/go/pkg/self_service/internal/application"
	"github.com/asgardeo/go/pkg/self_service/internal/association"
	"github.com/asgardeo/go/pkg/self_service/internal/authorized_apps"
	"github.com/asgardeo/go/pkg/self_service/internal/backup_code"
	"github.com/asgardeo/go/pkg/self_service/internal/consent"
	"github.com/asgardeo/go/pkg/self_service/internal/mfa"
	"github.com/asgardeo/go/pkg/self_service/internal/push"
	"github.com/asgardeo/go/pkg/self_service/internal/session"
)

type SessionListResponseModel = session.Sessions

type SessionModel = session.Session

type EnabledAuthenticatorsResponseModel = mfa.EnabledAuthenticators

type RemainingBackupCodesResponseModel = backup_code.RemainingBackupCodeResponse

type PushDeviceModel = push.DeviceDTO

type PasskeyModel = map[string]interface{}

type AssociatedUserModel = association.User

type FederatedAssociationModel = association.FederatedAssociation

type ApplicationListParamsModel = application.GetApplicationsParams

type ApplicationListResponseModel = application.ApplicationListResponse

type ApplicationModel = application.ApplicationResponse

type AuthorizedAppModel = authorized_apps.AuthorizedAppDTO

type ConsentListParamsModel = consent.ListConsentsOfLoggedInUserParams

type ConsentStateModel = consent.ListConsentsOfLoggedInUserParamsState

type ConsentSummaryModel = consent.ConsentSummary

type ConsentPurposeSummaryModel = consent.ConsentPurposeSummary

type ConsentModel = consent.ConsentResponse

type ConsentPurposeModel = consent.ConsentPurposeResponse
