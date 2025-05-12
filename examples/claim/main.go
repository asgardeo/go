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

package main

import (
	"context"
	"log"
	"time"

	"github.com/asgardeo/go/examples/common"
	"github.com/asgardeo/go/pkg/claim"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/sdk"
)

func main() {
	cfg := config.DefaultClientConfig().
		WithBaseURL("https://api.asgardeo.io/t/<tenant-domain>").
		WithTimeout(10*time.Second).
		WithClientCredentials(
			"client_id",
			"client_secret",
		)

	client, err := sdk.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	ctx := context.Background()

	// List Local Claims.
	claims, err := client.Claim.ListLocalClaims(ctx, nil)
	if err != nil {
		log.Printf("Error listing local claims: %v", err)
	} else {
		log.Printf("Found %d local claims.\n", len(*claims))
		log.Printf("Local claims: %v\n", common.ToJSONString(claims))
	}

	// List External Claims.
	externalClaims, err := client.Claim.ListExternalClaims(ctx, claim.ClaimDialectIDs.OIDC, nil)
	if err != nil {
		log.Printf("Error listing external claims: %v", err)
	} else {
		log.Printf("Found %d external claims.\n", len(*externalClaims))
		log.Printf("External claims: %v\n", common.ToJSONString(externalClaims))
	}

	// List OIDC Claims.
	oidcClaims, err := client.Claim.ListOIDCClaims(ctx, nil)
	if err != nil {
		log.Printf("Error listing OIDC claims: %v", err)
	} else {
		log.Printf("Found %d OIDC claims.\n", len(*oidcClaims))
		log.Printf("OIDC claims: %v\n", common.ToJSONString(oidcClaims))
	}
}
