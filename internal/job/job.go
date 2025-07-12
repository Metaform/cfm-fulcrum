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
	"github.com/metaform/cfm-fulcrum/internal/client"
	"time"
)

// JobHandler processes jobs from the Fulcrum Core job queue
type JobHandler struct {
	client client.FulcrumClient
	stats  struct {
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
func NewJobHandler(client client.FulcrumClient) *JobHandler {
	return &JobHandler{
		client: client,
	}
}

// PollAndProcessJobs polls for pending jobs and processes them
func (h *JobHandler) PollAndProcessJobs() error {
	// Get pending jobs
	jobs, err := h.client.GetPendingJobs()
	if err != nil {
		return fmt.Errorf("failed to get pending jobs: %w", err)
	}

	if len(jobs) == 0 {
		// log.Printf("Pending jobs not found")
		return nil
	}
	// First
	job := jobs[0]
	// Increment processed count
	h.stats.processed++
	// Claim the job
	if err := h.client.ClaimJob(job.ID); err != nil {
		// log.Printf("Failed to claim job %s: %v", job.ID, err)
		h.stats.failed++
		return err
	}
	//log.Printf("Processing job %s of type %s", job.ID, job.Action)
	// Process the job
	resp, err := h.processJob(job)
	if err != nil {
		// Mark job as failed
		//	log.Printf("Job %s failed: %v", job.ID, err)
		h.stats.failed++

		if failErr := h.client.FailJob(job.ID, err.Error()); failErr != nil {
			//	log.Printf("Failed to mark job %s as failed: %v", job.ID, failErr)
			return failErr
		}
	} else {
		// Job succeeded
		if complErr := h.client.CompleteJob(job.ID, resp); complErr != nil {
			//	log.Printf("Failed to mark job %s as completed: %v", job.ID, complErr)
			return complErr
		}
		h.stats.succeeded++
		//	log.Printf("Job %s completed successfully", job.ID)
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
	return nil, nil
}

// GetStats returns the job processing statistics
func (h *JobHandler) GetStats() (processed, succeeded, failed int) {
	return h.stats.processed, h.stats.succeeded, h.stats.failed
}
