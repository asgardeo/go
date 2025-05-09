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

import (
    "fmt"
    "net/url"
    "strings"
)

// extractOrigins gets allowed origins from a redirect URL
func extractOrigins(redirectURL string) ([]string, error) {
    if redirectURL == "" {
        return []string{}, nil
    }

    parsedURL, err := url.Parse(redirectURL)
    if err != nil {
        return nil, fmt.Errorf("authorized redirect URL is not in valid URL format: %w", err)
    }

    origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
    return []string{origin}, nil
}

// extractApplicationID extracts the application ID from a Location header URL
func extractApplicationID(locationHeader string) (string, error) {
    parsedURL, err := url.Parse(locationHeader)
    if err != nil {
        return "", fmt.Errorf("failed to parse Location header URL: %w", err)
    }
    
    return splitPath(parsedURL.Path), nil
}

// splitPath extracts the last segment from a path string
func splitPath(path string) string {
    parts := strings.Split(path, "/")
    
    for i := len(parts) - 1; i >= 0; i-- {
        if parts[i] != "" {
            return parts[i]
        }
    }
    
    return ""
}

// Helper functions for creating pointers to primitive types
func boolPtr(b bool) *bool {
    return &b
}

func stringPtr(s string) *string {
    return &s
}

func int64Ptr(i int64) *int64 {
    return &i
}

func intPtr(i int) *int {
	return &i
}
