package scenario

import (
	"fmt"
	"github.com/metaform/cfm-fulcrum/internal/client"
)

type ServiceRequest struct {
	Name          string                 `json:"name"`
	Properties    map[string]interface{} `json:"properties"`
	AgentTags     []string               `json:"agentTags"`
	AgentID       string                 `json:"agentId"`
	ServiceTypeID string                 `json:"serviceTypeId"`
	GroupID       string                 `json:"groupId"`
}

func RunCreateTenantDeploymentCommand(cfg Config) error {
	fmt.Println("Starting tenant deployment...")

	apiClient := client.NewApiClient(pmanagerBaseUrl, fulcrumCoreBaseUrl, cfmAgentBaseUrl)

	serviceRequest := ServiceRequest{
		AgentID:       cfg.AgentID,
		Name:          "tenant-deployment",
		Properties:    make(map[string]interface{}),
		AgentTags:     []string{},
		ServiceTypeID: cfmTenantServiceType,
		GroupID:       cfg.ServiceGroupID,
	}

	serviceRequest.Properties["tenantDid"] = "did:web:tenant.example.com"
	_, err := apiClient.PostToFulcrumCore("services", adminToken, serviceRequest)
	if err != nil {
		return err
	}

	fmt.Println("Service created")
	return nil
}
