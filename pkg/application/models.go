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

type LoginFlowGenerateResponseModel = internal.LoginFlowGenerateResponse

type LoginFlowStatusResponseModel = internal.LoginFlowStatusResponse

type LoginFlowResultResponseModel = internal.LoginFlowResultResponse

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
