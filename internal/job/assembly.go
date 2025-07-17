package job

import (
	"github.com/metaform/cfm-fulcrum/internal/client"
	"github.com/metaform/connector-fabric-manager/common/system"
	"time"
)

type JobServiceAssembly struct {
	system.DefaultServiceAssembly
	handler     *JobHandler
	stopChannel chan struct{}
}

func (a *JobServiceAssembly) Name() string {
	return "Job Service"
}

func (d *JobServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{client.FulcrumClientKey}
}

func (a *JobServiceAssembly) Init(context *system.InitContext) error {
	fulcrumClient := context.Registry.Resolve(client.FulcrumClientKey).(client.FulcrumClient)
	apiClient := context.Registry.Resolve(client.ApiClientKey).(client.ApiClient)

	a.handler = NewJobHandler(fulcrumClient, apiClient, context.LogMonitor)
	return nil
}

func (a *JobServiceAssembly) Start(ctx *system.StartContext) error {
	if a.handler == nil {
		return nil
	}
	a.stopChannel = make(chan struct{})

	go func() {
		ticker := time.NewTicker(30 * time.Second) // TODO configure
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx.LogMonitor.Infof("Polling jobs")
				if err := a.handler.PollAndProcessJobs(); err != nil {
					ctx.LogMonitor.Infof("Error polling jobs: %w", err)
				}
			case <-a.stopChannel:
				ctx.LogMonitor.Infof("Stopping job service")
				return
			}
		}
	}()
	return nil
}

func (a *JobServiceAssembly) Finalize() error {
	if a.handler == nil || a.stopChannel == nil {
		return nil
	}
	a.stopChannel <- struct{}{}
	return nil
}
