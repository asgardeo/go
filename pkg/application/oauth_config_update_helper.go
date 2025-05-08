package application

// NewOAuthConfigUpdate creates a new ApplicationOAuthConfigUpdateModel with default values
func NewOAuthConfigUpdate() *ApplicationOAuthConfigUpdateModel {
	return &ApplicationOAuthConfigUpdateModel{}
}

// WithAccessTokenAttributes sets access token attributes
func (c *ApplicationOAuthConfigUpdateModel) WithAccessTokenAttributes(attrs []string) *ApplicationOAuthConfigUpdateModel {
	c.AccessTokenAttributes = &attrs
	return c
}

// WithApplicationAccessTokenExpiry sets application access token expiry in seconds
func (c *ApplicationOAuthConfigUpdateModel) WithApplicationAccessTokenExpiry(seconds int64) *ApplicationOAuthConfigUpdateModel {
	c.ApplicationAccessTokenExpiryInSeconds = &seconds
	return c
}

// WithUserAccessTokenExpiry sets user access token expiry in seconds
func (c *ApplicationOAuthConfigUpdateModel) WithUserAccessTokenExpiry(seconds int64) *ApplicationOAuthConfigUpdateModel {
	c.UserAccessTokenExpiryInSeconds = &seconds
	return c
}

// WithAllowedOrigins sets allowed origins for CORS
func (c *ApplicationOAuthConfigUpdateModel) WithAllowedOrigins(origins []string) *ApplicationOAuthConfigUpdateModel {
	c.AllowedOrigins = &origins
	return c
}

// WithCallbackURLs sets callback URLs
func (c *ApplicationOAuthConfigUpdateModel) WithCallbackURLs(urls []string) *ApplicationOAuthConfigUpdateModel {
	c.CallbackURLs = &urls
	return c
}

// WithRefreshTokenExpiry sets refresh token expiry in seconds
func (c *ApplicationOAuthConfigUpdateModel) WithRefreshTokenExpiry(seconds int64) *ApplicationOAuthConfigUpdateModel {
	c.RefreshTokenExpiryInSeconds = &seconds
	return c
}
