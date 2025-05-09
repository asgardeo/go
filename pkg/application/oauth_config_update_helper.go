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
