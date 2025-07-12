package job

import (
	"github.com/metaform/cfm-fulcrum/internal/client"
	"github.com/metaform/connector-fabric-manager/common/system"
)

type JobServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (a *JobServiceAssembly) Name() string {
	return "Job Service"
}

func (d *JobServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{client.ClientKey}
}

func (a *JobServiceAssembly) Init(context *system.InitContext) error {
	fulcrumClient := context.Registry.Resolve(client.ClientKey).(client.FulcrumClient)
	NewJobHandler(fulcrumClient)
	return nil
}
