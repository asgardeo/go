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

package main

import (
	"context"
	"log"
	"time"

	"github.com/asgardeo/go/examples/common"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/sdk"
)

func main() {

	// Create a configuration with the access token of the user, with the internal_login scope
	cfg := config.DefaultClientConfig().
		WithBaseURL("https://api.asgardeo.io/t/<tenant-domain>").
		WithTimeout(10 * time.Second).
		WithToken("<user_access_token>")

	client, err := sdk.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// To use a different token per request, configure WithContextToken() instead of WithToken(),
	// and set the token on the context: ctx = config.ContextWithToken(ctx, "<user_access_token>")
	ctx := context.Background()

	// Get the authenticated user.
	me, err := client.User.GetMe(ctx, "userName,name")
	if err != nil {
		log.Printf("Error getting the user: %v", err)
	} else {
		log.Printf("Found user: %s\n", common.ToJSONString(me))
	}

	// List sessions.
	sessions, err := client.SelfService.ListSessions(ctx)
	if err != nil {
		log.Printf("Error listing sessions: %v", err)
	} else {
		log.Printf("Sessions: %s\n", common.ToJSONString(sessions))
	}

	// Terminate a session by ID.
	err = client.SelfService.TerminateSession(ctx, "session_id")
	if err != nil {
		log.Printf("Error terminating session: %v", err)
	} else {
		log.Printf("Terminated session.\n")
	}
}
