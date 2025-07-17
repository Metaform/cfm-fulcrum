//  Copyright (c) 2025 Metaform Systems, Inc
//
//  This program and the accompanying materials are made available under the
//  terms of the Apache License, Version 2.0 which is available at
//  https://www.apache.org/licenses/LICENSE-2.0
//
//  SPDX-License-Identifier: Apache-2.0
//
//  Contributors:
//       Metaform Systems, Inc. - initial API and implementation
//

package scenario

import (
	"encoding/json"
	"fmt"
	"github.com/metaform/cfm-fulcrum/internal/client"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"time"
)

// RunOnboardCommand executes the complete onboarding process and returns the config
func RunOnboardCommand() (*Config, error) {
	fmt.Println("Starting onboarding process...")

	apiClient := client.NewApiClient(pmanagerBaseUrl, fulcrumCoreBaseUrl, cfmAgentBaseUrl)

	err := CreateTestActivityDefinition(apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create test activity definition: %w", err)
	}

	err = CreateTestDeploymentDefinition(apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create test deployment definition: %w", err)
	}

	providerId, err := CreateProvider("Test Provider", apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	fmt.Printf("Created Fulcrum provider: %s\n", providerId)

	serviceGroupId, err := CreateServiceGroup("CFM Service Group", providerId, apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create service group: %w", err)
	}

	agentId, err := CreateAgent("Test Agent", providerId, apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}
	fmt.Printf("Created Fulcrum agent: %s\n", agentId)

	agentToken, err := CreateAgentToken("Test Agent Token", agentId, apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent token: %w", err)
	}

	fmt.Println("Created agent token")

	err = UpdateToken(agentToken, apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	fmt.Println("Updated CFM agent with token")
	fmt.Println("Onboarding process completed successfully")

	return &Config{
		ProviderId:     providerId,
		AgentID:        agentId,
		ServiceGroupID: serviceGroupId,
	}, nil
}

func CreateTestActivityDefinition(apiClient *client.ApiClient) error {
	requestBody := api.ActivityDefinition{
		Type:        "test.activity",
		Description: "Performs a test activity",
	}

	return apiClient.PostToPManager("activity-definition", requestBody)
}

func CreateTestDeploymentDefinition(apiClient *client.ApiClient) error {
	requestBody := api.DeploymentDefinition{
		Type:       "test.deployment",
		ApiVersion: "v1",
		Resource: api.Resource{
			Group:       "deployments.example.com",
			Singular:    "TestDeployment",
			Plural:      "TestDeployments",
			Description: "Test deployment",
		},
		Versions: []api.Version{
			{
				Version: "1.0.0",
				Active:  true,
				Activities: []api.Activity{
					{
						ID:   "activity1",
						Type: "test-activity",
					},
				},
			},
		},
	}

	return apiClient.PostToPManager("deployment-definition", requestBody)
}

//func CreateTestDeployment() error {
//	requestBody := api.DeploymentManifest{
//		DeploymentType: "test.deployment",
//		ID:             uuid.New().String(),
//		Payload:        make(map[string]any),
//	}
//
//	client := client.NewApiClient()
//	return client.PostToPManager("deployment", requestBody)
//}

type CreateProviderRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func CreateProvider(name string, apiClient *client.ApiClient) (string, error) {
	requestBody := CreateProviderRequest{
		Name:   name,
		Status: "Enabled",
	}

	body, err := apiClient.PostToFulcrumCore("participants", adminToken, requestBody)
	if err != nil {
		return "", err
	}

	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response["id"], nil
}

type CreateAgentRequest struct {
	Name        string   `json:"name"`
	ProviderID  string   `json:"providerId"`
	AgentTypeID string   `json:"agentTypeId"`
	Tags        []string `json:"tags"`
}

// CreateAgent creates a new agent using the Fulcrum Core API
func CreateAgent(name string, providerID string, apiClient *client.ApiClient) (string, error) {
	requestBody := CreateAgentRequest{
		Name:        name,
		ProviderID:  providerID,
		AgentTypeID: seedAgentType,
		Tags:        []string{"cfm"},
	}

	body, err := apiClient.PostToFulcrumCore("agents", adminToken, requestBody)
	if err != nil {
		return "", err
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response["id"].(string), nil
}

type CreateAgentTokenRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ScopeID     string `json:"scopeId"`
	ExpiresAt   string `json:"expiresAt"`
	Role        string `json:"role"`
}

func CreateAgentToken(name string, agentID string, apiClient *client.ApiClient) (string, error) {
	requestBody := CreateAgentTokenRequest{
		Name:        name,
		Description: "Agent token",
		ScopeID:     agentID,
		ExpiresAt:   time.Now().AddDate(50, 0, 0).Format(time.RFC3339),
		Role:        "agent",
	}

	body, err := apiClient.PostToFulcrumCore("tokens", adminToken, requestBody)
	if err != nil {
		return "", err
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response["value"].(string), nil
}

type ServiceGroupRequest struct {
	Name       string `json:"name"`
	ConsumerID string `json:"consumerID"`
}

func CreateServiceGroup(name, consumerID string, apiClient *client.ApiClient) (string, error) {
	requestBody := ServiceGroupRequest{
		Name:       name,
		ConsumerID: consumerID,
	}

	body, err := apiClient.PostToFulcrumCore("service-groups", adminToken, requestBody)
	if err != nil {
		return "", err
	}

	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response["id"], nil
}

type UpdateTokenRequest struct {
	Token string `json:"token"`
}

// UpdateToken updates the token on the CFM agent
func UpdateToken(token string, apiClient *client.ApiClient) error {
	requestBody := UpdateTokenRequest{
		Token: token,
	}

	return apiClient.PostToCFMAgent("fulcrum-token", requestBody)
}
