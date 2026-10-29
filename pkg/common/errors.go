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

package common

import (
	"fmt"
	"slices"
)

// APIError is returned when the server responds with an unexpected status code
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("status %d, body: %s", e.StatusCode, string(e.Body))
}

// CheckStatus returns an *APIError unless the status code is one of the expected ones
func CheckStatus(action string, statusCode int, body []byte, expected ...int) error {
	if slices.Contains(expected, statusCode) {
		return nil
	}
	return fmt.Errorf("failed to %s: %w", action, &APIError{StatusCode: statusCode, Body: body})
}
