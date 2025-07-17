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

package job

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/metaform/cfm-fulcrum/internal/client"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"time"
)

// JobHandler processes jobs from the Fulcrum Core job queue
type JobHandler struct {
	fulcrumClient client.FulcrumClient
	apiClient     client.ApiClient
	monitor       monitor.LogMonitor
	stats         struct {
		processed int
		succeeded int
		failed    int
	}
}

// JobResources represents the resources in a job response
type JobResources struct {
	TS time.Time `json:"ts"`
}

// JobResponse represents the response for a job
type JobResponse struct {
	Resources  JobResources `json:"resources"`
	ExternalID *string      `json:"externalId"`
}

type VMProps struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
}

// NewJobHandler creates a new job handler
func NewJobHandler(fulcrumClient client.FulcrumClient, apiClient client.ApiClient, monitor monitor.LogMonitor) *JobHandler {
	return &JobHandler{
		fulcrumClient: fulcrumClient,
		apiClient:     apiClient,
		monitor:       monitor,
	}
}

// PollAndProcessJobs polls for pending jobs and processes them
func (h *JobHandler) PollAndProcessJobs() error {
	// Get pending jobs
	jobs, err := h.fulcrumClient.GetPendingJobs()
	if err != nil {
		return fmt.Errorf("failed to get pending jobs: %w", err)
	}

	if len(jobs) == 0 {
		h.monitor.Infof("Pending jobs not found")
		return nil
	}
	job := jobs[0]
	h.stats.processed++

	// Claim the job
	if err := h.fulcrumClient.ClaimJob(job.ID); err != nil {
		h.stats.failed++
		return err
	}

	// Process the job
	resp, err := h.processJob(job)
	if err != nil {
		// Mark job as failed
		h.stats.failed++
		if failErr := h.fulcrumClient.FailJob(job.ID, err.Error()); failErr != nil {
			//	log.Printf("Failed to mark job %s as failed: %v", job.ID, failErr)
			return failErr
		}
	} else {
		// Job succeeded
		if complErr := h.fulcrumClient.CompleteJob(job.ID, resp); complErr != nil {
			//	log.Printf("Failed to mark job %s as completed: %v", job.ID, complErr)
			return complErr
		}
		h.stats.succeeded++
	}

	return nil
}

// processJob processes a job based on its type
func (h *JobHandler) processJob(job *client.Job) (any, error) {
	switch job.Action {
	case client.JobActionServiceCreate:
	case client.JobActionServiceColdUpdate, client.JobActionServiceHotUpdate:
	case client.JobActionServiceStart:
	case client.JobActionServiceStop:
	case client.JobActionServiceDelete:
	default:
		return nil, fmt.Errorf("unknown job type: %s", job.Action)
	}

	requestBody := api.DeploymentManifest{
		DeploymentType: "test.deployment",
		ID:             uuid.New().String(),
		Payload:        make(map[string]any),
	}

	fmt.Printf("Processing job %s of type %s", job.ID, job.Action)
	err := h.apiClient.PostToPManager("deployment", requestBody)
	if err != nil {
		h.monitor.Severef("**********error in job handler **********: %w", err)
		return nil, err
	}
	return nil, nil
}

// GetStats returns the job processing statistics
func (h *JobHandler) GetStats() (processed, succeeded, failed int) {
	return h.stats.processed, h.stats.succeeded, h.stats.failed
}
