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

package user

import "github.com/asgardeo/go/pkg/user/internal"

type UserCreateModel struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// convertToUserCreateModel converts the UserCreateModel to the internal.AddUserJSONBody model.
func convertToAddUserJSONBodyModel(user UserCreateModel) internal.AddUserJSONBody {
	return internal.AddUserJSONBody{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Name: internal.Name{
			GivenName:  user.FirstName,
			FamilyName: user.LastName,
		},
		Emails: []internal.Email{
			{
				Primary: true,
				Value:   user.Email,
			},
		},
	}
}
