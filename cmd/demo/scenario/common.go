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

const (
	pmanagerBaseUrl      = "http://localhost:8181"
	fulcrumCoreBaseUrl   = "http://localhost:8080/api/v1"
	cfmAgentBaseUrl      = "http://localhost:8383"
	adminToken           = "admin-test-token"
	seedAgentType        = "f47ac10b-58cc-4372-a567-0e02b2c3d479" // created from Fulcrum seeding
	cfmTenantServiceType = "01940a2e-7b8f-7c4d-9e5a-3f2b1c8d9e0f" // created from Fulcrum seeding
)

type Config struct {
	ProviderId     string `json:"providerId"`
	AgentID        string `json:"agentId"`
	ServiceGroupID string `json:"serviceGroupId"`
}
