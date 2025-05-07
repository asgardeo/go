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
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/asgardeo/go/pkg/application"
	"github.com/asgardeo/go/pkg/config"
	"github.com/asgardeo/go/pkg/sdk"
)

func main() {

	// Create a configuration with a client credentials grant type
	cfg := config.DefaultClientConfig().
		WithBaseURL("https://api.asgardeo.io/t/<tenant-domain>").
		WithTimeout(10*time.Second).
		WithClientCredentials(
			"client_id",
			"client_secret",
		)

	// Create a client with the client credentials configuration
	client, err := sdk.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create SDK client: %v", err)
	}

	// Use the client with token authentication
	ctx := context.Background()

	// List applications.
	apps, err := client.Application.List(ctx, 10, 0)
	if err != nil {
		log.Printf("Error listing users: %v", err)
	} else {
		fmt.Printf("Found %d applications\n", len(*apps.Applications))
	}

	// Create a SPA with name and redirect URL
	spa, err := client.Application.CreateSinglePageApp(
		context.Background(),
		"spa_name",
		"https://example.com/callback",
	)

	if err != nil {
		log.Printf("Error creating application: %v", err)
	}

	fmt.Printf("Created SPA:\n%s\n", toJSONString(spa))

	// Create a mobile app with name and redirect URL
	mobileApp, err := client.Application.CreateMobileApp(
		context.Background(),
		"mobile_name",
		"https://example.com/callback",
	)

	if err != nil {
		log.Printf("Error creating application: %v", err)
	}

	fmt.Printf("Created Mobile App with %s \n", toJSONString(mobileApp))

	// Create a M2M app with name
	m2mApp, err := client.Application.CreateM2MApp(
		context.Background(),
		"m2m_name",
	)

	if err != nil {
		log.Printf("Error creating application: %v", err)
	}

	fmt.Printf("Created M2M App with %s \n", toJSONString(m2mApp))

	// Get application by name
	app, err := client.Application.GetByName(context.Background(), "app_name")
	if err != nil {
		fmt.Printf("Error retrieiving application: %v\n", err)
		return
	}
	fmt.Printf("Found app %s with %s \n", app.Name, app)

	// Get application by client ID
	app, err = client.Application.GetByClienId(context.Background(), "app_client_id")
	if err != nil {
		fmt.Printf("Error retrieving application: %v\n", err)
		return
	}
	fmt.Printf("Found app %s with %s \n", app.Name, app)

	// Update application basic info.
	updating_app_name := "updating_app_name"
	updating_app_description := "updating_app_description"

	updatingApplication := application.ApplicationBasicInfoUpdateModel{
		Name:        &updating_app_name,
		Description: &updating_app_description,
	}

	updatedApp, err := client.Application.UpdateBasicInfo(ctx, "updating_app_id", updatingApplication)
	if err != nil {
		fmt.Printf("Error updating application: %v\n", err)
		return
	}
	fmt.Printf("Updated app %s with %s \n", updatedApp.Name, updatedApp)

	// Authorize API.
	id := "api_resource_uuid"
	policyIdentifier := "RBAC"
	scopes := []string{"scope1", "scope2"}
	authorizedAPI := application.AuthorizedAPICreateModel{
		Id:               &id,
		PolicyIdentifier: &policyIdentifier,
		Scopes:           &scopes,
	}
	err = client.Application.AuthorizeAPI(ctx, "app_uuid", authorizedAPI)
	if err != nil {
		log.Printf("Error authorizing API: %v", err)
	} else {
		log.Printf("API authorized successfully.")
	}

	// Generate a login flow.
	prompt := "Username and password as the first step and email OTP as the second step."
	loginFlowResponse, err := client.Application.GenerateLoginFlow(ctx, prompt)
	fmt.Printf("Login flow response: %+v\n", loginFlowResponse)
	if err != nil {
		log.Printf("Error generating login flow: %v", err)
		return
	}
	log.Printf("Login flow initiated. Flow ID: %s", *loginFlowResponse.OperationId)

	// Poll for the login flow generation status.
	flowId := loginFlowResponse.OperationId
	var statusResponse *application.LoginFlowStatusResponseModel
	for {
		statusResponse, err = client.Application.GetLoginFlowGenerationStatus(ctx, *flowId)
		if err != nil {
			log.Printf("Error getting login flow generation status: %v", err)
			return
		}
		if statusResponse.Status != nil {
			allTrue := true
			for _, v := range *statusResponse.Status {
				if v != true {
					allTrue = false
					break
				}
			}
			if allTrue {
				log.Printf("Login flow generation completed successfully.")
				break
			}
		}
		// If the status is not complete, wait and poll again.
		log.Printf("Login flow generation in progress. Retrying...")
		time.Sleep(2 * time.Second)
	}

	// Retrieve the login flow generation result.
	resultResponse, err := client.Application.GetLoginFlowGenerationResult(ctx, *flowId)
	if err != nil {
		log.Printf("Error getting login flow generation result: %v", err)
		return
	}
	log.Printf("Login flow generation result: %+v", resultResponse.Data)
}

func toJSONString(app interface{}) string {
	jsonData, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		log.Printf("Error marshaling app object: %v", err)
		return ""
	}
	return string(jsonData)
}
