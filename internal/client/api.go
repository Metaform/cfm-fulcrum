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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiClient struct {
	pmanagerBaseUrl    string
	fulcrumCoreBaseUrl string
	cfmAgentBaseUrl    string
	client             http.Client
}

func NewApiClient(pmanagerBaseUrl string, fulcrumCoreBaseUrl string, cfmAgentBaseUrl string) *ApiClient {
	return &ApiClient{
		pmanagerBaseUrl:    pmanagerBaseUrl,
		fulcrumCoreBaseUrl: fulcrumCoreBaseUrl,
		cfmAgentBaseUrl:    cfmAgentBaseUrl,
		client:             http.Client{},
	}
}

// PostToFulcrumCore makes a POST request to Fulcrum Core API with authentication
func (c *ApiClient) PostToFulcrumCore(endpoint string, token string, payload any) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.fulcrumCoreBaseUrl, endpoint)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	return c.postRequest(url, payload, headers)
}

// PostToPManager makes a POST request to Process Manager API
func (c *ApiClient) PostToPManager(endpoint string, payload any) error {
	url := fmt.Sprintf("%s/%s", c.pmanagerBaseUrl, endpoint)
	_, err := c.postRequest(url, payload, nil)
	return err
}

// PostToCFMAgent makes a POST request to CFM Agent API
func (c *ApiClient) PostToCFMAgent(endpoint string, payload any) error {
	url := fmt.Sprintf("%s/%s", c.cfmAgentBaseUrl, endpoint)
	_, err := c.postRequest(url, payload, nil)
	return err
}

// postRequest handles POST requests with JSON payload
func (c *ApiClient) postRequest(url string, payload any, headers map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
