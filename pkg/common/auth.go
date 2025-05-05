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

package common

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asgardeo/go/pkg/config"
)

// CreateAuthRequestEditorFunc returns a function that adds authentication to requests
func CreateAuthRequestEditorFunc(cfg *config.ClientConfig) interface{} {
	// Return a function that matches the RequestEditorFn signature
	// The caller will need to cast this to the appropriate type
	return func(ctx context.Context, req *http.Request) error {

		token, err := cfg.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get authentication token: %w", err)
		}

		if token != "" {
			// Add Authorization header with Bearer token
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		}
		return nil
	}
}
