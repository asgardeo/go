package application

// NewBasicInfoUpdate creates a new ApplicationBasicInfoUpdateModel with default values
func NewBasicInfoUpdate() *ApplicationBasicInfoUpdateModel {
	return &ApplicationBasicInfoUpdateModel{}
}

// WithName sets application name
func (c *ApplicationBasicInfoUpdateModel) WithName(name string) *ApplicationBasicInfoUpdateModel {
	c.Name = &name
	return c
}

// WithDescription sets application description
func (c *ApplicationBasicInfoUpdateModel) WithDescription(description string) *ApplicationBasicInfoUpdateModel {
	c.Description = &description
	return c
}

// WithImageUrl sets application image URL
func (c *ApplicationBasicInfoUpdateModel) WithImageUrl(imageUrl string) *ApplicationBasicInfoUpdateModel {
	c.ImageUrl = &imageUrl
	return c
}

// WithAccessUrl sets application access URL
func (c *ApplicationBasicInfoUpdateModel) WithAccessUrl(accessUrl string) *ApplicationBasicInfoUpdateModel {
	c.AccessUrl = &accessUrl
	return c
}

// WithLogoutReturnUrl sets application logout return URL
func (c *ApplicationBasicInfoUpdateModel) WithLogoutReturnUrl(logoutReturnUrl string) *ApplicationBasicInfoUpdateModel {
	c.LogoutReturnUrl = &logoutReturnUrl
	return c
}
