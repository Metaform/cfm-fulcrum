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

//go:generate mockery --name FulcrumClient --filename job_mock.go --with-expecter --outpkg mocks --dir . --output ./mocks

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync/atomic"
	"time"
)

// JobAction represents the type of job
type JobAction string

const (
	JobActionServiceCreate     JobAction = "ServiceCreate"
	JobActionServiceStart      JobAction = "ServiceStart"
	JobActionServiceStop       JobAction = "ServiceStop"
	JobActionServiceHotUpdate  JobAction = "ServiceHotUpdate"
	JobActionServiceColdUpdate JobAction = "ServiceColdUpdate"
	JobActionServiceDelete     JobAction = "ServiceDelete"
)

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusPending    JobStatus = "Pending"
	JobStatusProcessing JobStatus = "Processing"
	JobStatusCompleted  JobStatus = "Completed"
	JobStatusFailed     JobStatus = "Failed"
)

// Job represents a job from the Fulcrum Core job queue
type Job struct {
	ID       string    `json:"id"`
	Action   JobAction `json:"action"`
	Status   JobStatus `json:"status"`
	Priority int       `json:"priority"`
	Service  struct {
		ID                string  `json:"id"`
		Name              string  `json:"name"`
		ExternalID        *string `json:"externalId"`
		CurrentProperties *struct {
			CPU    int `json:"cpu"`
			Memory int `json:"memory"`
		} `json:"currentProperties"`
		TargetProperties *struct {
			CPU    int `json:"cpu"`
			Memory int `json:"memory"`
		} `json:"targetProperties"`
	} `json:"service"`
}

// MetricEntry represents a single metric measurement
type MetricEntry struct {
	ExternalID string  `json:"externalId"`
	ResourceID string  `json:"resourceId"`
	Value      float64 `json:"value"`
	TypeName   string  `json:"typeName"`
}

// FulcrumClient defines the interface for communication with the Fulcrum Core API
type FulcrumClient interface {
	UpdateAgentStatus(status string) error
	GetAgentInfo() (map[string]any, error)
	GetPendingJobs() ([]*Job, error)
	ClaimJob(jobID string) error
	CompleteJob(jobID string, resources any) error
	FailJob(jobID string, errorMessage string) error
	ReportMetric(metrics *MetricEntry) error

	UpdateToken(token string) error
}

type HTTPFulcrumClient struct {
	baseURL    string
	httpClient *http.Client
	token      atomic.Pointer[string] // Agent authentication token
}

func NewHTTPFulcrumClient(baseURL string, token string) FulcrumClient {
	client := &HTTPFulcrumClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	client.token.Store(&token)
	return client
}

func (c *HTTPFulcrumClient) UpdateToken(token string) error {
	c.token.Store(&token)
	return nil
}

// UpdateAgentStatus updates the agent's status in Fulcrum Core
func (c *HTTPFulcrumClient) UpdateAgentStatus(status string) error {
	reqBody, err := json.Marshal(map[string]any{
		"status": status,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal status update request: %w", err)
	}

	resp, err := c.put("/api/v1/agents/me/status", reqBody)
	if err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update agent status, status: %d", resp.StatusCode)
	}

	return nil
}

// GetAgentInfo retrieves the agent's information from Fulcrum Core
func (c *HTTPFulcrumClient) GetAgentInfo() (map[string]any, error) {
	resp, err := c.get("/api/v1/agents/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get agent info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get agent info, status: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode agent info response: %w", err)
	}

	return result, nil
}

// GetPendingJobs retrieves pending jobs for this agent
func (c *HTTPFulcrumClient) GetPendingJobs() ([]*Job, error) {
	resp, err := c.get("/api/v1/jobs/pending")
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pending jobs, status: %d", resp.StatusCode)
	}

	var jobs []*Job

	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("failed to decode jobs response: %w", err)
	}

	return jobs, nil
}

// ClaimJob claims a job for processing
func (c *HTTPFulcrumClient) ClaimJob(jobID string) error {
	resp, err := c.post(fmt.Sprintf("/api/v1/jobs/%s/claim", jobID), nil)
	if err != nil {
		return fmt.Errorf("failed to claim job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to claim job, status: %d", resp.StatusCode)
	}

	return nil
}

// CompleteJob marks a job as completed with results
func (c *HTTPFulcrumClient) CompleteJob(jobID string, response any) error {
	reqBody, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal job completion request: %w", err)
	}

	resp, err := c.post(fmt.Sprintf("/api/v1/jobs/%s/complete", jobID), reqBody)
	if err != nil {
		return fmt.Errorf("failed to complete job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to complete job, status: %d", resp.StatusCode)
	}

	return nil
}

// FailJob marks a job as failed with an error message
func (c *HTTPFulcrumClient) FailJob(jobID string, errorMessage string) error {
	reqBody, err := json.Marshal(map[string]any{
		"errorMessage": errorMessage,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal job failure request: %w", err)
	}

	resp, err := c.post(fmt.Sprintf("/api/v1/jobs/%s/fail", jobID), reqBody)
	if err != nil {
		return fmt.Errorf("failed to mark job as failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to mark job as failed, status: %d", resp.StatusCode)
	}

	return nil
}

// ReportMetrics sends collected metrics to Fulcrum Core
func (c *HTTPFulcrumClient) ReportMetric(metric *MetricEntry) error {
	reqBody, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics request: %w", err)
	}

	resp, err := c.post("/api/v1/metric-entries", reqBody)
	if err != nil {
		return fmt.Errorf("failed to report metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to report metrics, status: %d", resp.StatusCode)
	}

	return nil
}

// Helper methods for HTTP requests
func (c *HTTPFulcrumClient) get(endpoint string) (*http.Response, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("******" + c.getToken())
	req.Header.Set("Authorization", "Bearer "+c.getToken())
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *HTTPFulcrumClient) post(endpoint string, body []byte) (*http.Response, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.getToken())

	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *HTTPFulcrumClient) put(endpoint string, body []byte) (*http.Response, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.getToken())
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *HTTPFulcrumClient) getToken() string {
	tokenPtr := c.token.Load()
	if tokenPtr == nil {
		return ""
	}
	return *tokenPtr
}
