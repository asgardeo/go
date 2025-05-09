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
