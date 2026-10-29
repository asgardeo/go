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
	"github.com/asgardeo/go/pkg/consent"
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

	ctx := context.Background()

	// List the active consents of the user.
	username := "<username>"
	state := "ACTIVE"
	receipts, err := client.Consent.List(ctx, &consent.ConsentListParamsModel{
		PiiPrincipalId: &username,
		State:          &state,
	})
	if err != nil {
		log.Printf("Error listing consents: %v", err)
	} else {
		log.Printf("Found %d consents.\n", len(receipts))
		log.Printf("Consents: %s\n", common.ToJSONString(receipts))
	}

	// Get a consent receipt by ID.
	receipt, err := client.Consent.GetReceipt(ctx, "consent_receipt_id")
	if err != nil {
		log.Printf("Error getting consent receipt: %v", err)
	} else {
		log.Printf("Found consent receipt: %s\n", common.ToJSONString(receipt))
	}

	// Revoke a consent receipt by ID.
	err = client.Consent.RevokeReceipt(ctx, "consent_receipt_id")
	if err != nil {
		log.Printf("Error revoking consent receipt: %v", err)
	} else {
		log.Printf("Revoked consent receipt.\n")
	}
}
