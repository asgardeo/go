package application

import (
	"github.com/asgardeo/go/pkg/application/internal"
)

type ApplicationBasicInfoResponseModel struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ClientId         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	RedirectURL      string `json:"redirect_url"`
	AuthorizedScopes string `json:"scope"`
}

type ApplicationListResponseModel = internal.ApplicationListResponse

type AuthorizedAPICreateModel = internal.AddAuthorizedAPIJSONRequestBody

// ApplicationBasicInfoUpdateModel defines a simplified model for updating basic application information
type ApplicationBasicInfoUpdateModel struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	ImageUrl        *string `json:"imageUrl,omitempty"`
	AccessUrl       *string `json:"accessUrl,omitempty"`
	LogoutReturnUrl *string `json:"logoutReturnUrl,omitempty"`
}

// ApplicationOAuthConfigUpdateModel contains only the fields that can be updated in OAuth configuration
type ApplicationOAuthConfigUpdateModel struct {
	// Access token fields
	AccessTokenAttributes                 *[]string `json:"accessTokenAttributes,omitempty"`
	ApplicationAccessTokenExpiryInSeconds *int64    `json:"applicationAccessTokenExpiryInSeconds,omitempty"`
	UserAccessTokenExpiryInSeconds        *int64    `json:"userAccessTokenExpiryInSeconds,omitempty"`

	// CORS and redirect URIs
	AllowedOrigins *[]string `json:"allowedOrigins,omitempty"`
	CallbackURLs   *[]string `json:"callbackURLs,omitempty"`

	// Logout configuration
	Logout *internal.OIDCLogoutConfiguration `json:"logout,omitempty"`

	// Refresh token config
	RefreshTokenExpiryInSeconds *int64 `json:"refreshTokenExpiryInSeconds,omitempty"`
}

type LoginFlowGenerateResponseModel = internal.LoginFlowGenerateResponse

type LoginFlowStatusResponseModel = internal.LoginFlowStatusResponse

type LoginFlowResultResponseModel struct {
	Data   *LoginFlowUpdateModel `json:"data,omitempty"`
	Status *internal.StatusEnum  `json:"status,omitempty"`
}

type LoginFlowUpdateModel = internal.AuthenticationSequence

type LoginFlowStepModel = internal.AuthenticationStepModel

type AuthenticatorModel = internal.Authenticator

type LoginFlowTypeModel = internal.AuthenticationSequenceType

// convertToApplicationPatchModel converts the public ApplicationBasicInfoUpdateModel to the internal PatchApplicationJSONRequestBody
func convertToApplicationPatchModel(model ApplicationBasicInfoUpdateModel) internal.PatchApplicationJSONRequestBody {
	return internal.PatchApplicationJSONRequestBody{
		Name:            model.Name,
		Description:     model.Description,
		ImageUrl:        model.ImageUrl,
		AccessUrl:       model.AccessUrl,
		LogoutReturnUrl: model.LogoutReturnUrl,
	}
}

func convertToLoginFlowResultResponseModel(model internal.LoginFlowResultResponse) LoginFlowResultResponseModel {
	loginFlowUpdateData := convertToLoginFlowUpdateModel(model.Data)
	return LoginFlowResultResponseModel{
		Data:   &loginFlowUpdateData,
		Status: model.Status,
	}
}

func convertToLoginFlowUpdateModel(data *map[string]interface{}) LoginFlowUpdateModel {
	var loginFlowUpdate LoginFlowUpdateModel
	if data != nil {
		loginFlowUpdate = LoginFlowUpdateModel{
			AttributeStepId: (*data)["attributeStepId"].(*int),
			Steps:           (*data)["steps"].(*[]LoginFlowStepModel),
			SubjectStepId:   (*data)["subjectStepId"].(*int),
			Type:            (*data)["type"].(*LoginFlowTypeModel),
		}
	}
	return loginFlowUpdate
}
