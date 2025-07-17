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

package launcher

import (
	"fmt"
	"github.com/metaform/cfm-fulcrum/internal/client"
	"github.com/metaform/cfm-fulcrum/internal/job"
	"github.com/metaform/cfm-fulcrum/internal/management"
	"github.com/metaform/cfm-fulcrum/internal/sysconfig"
	"github.com/metaform/connector-fabric-manager/assembly/httpclient"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	configPrefix = "cfm-agent"
	agentName    = "Fulcrum CFM Agent"
	defaultPort  = 8080
	httpKey      = "httpPort"
)

func LaunchAndWaitSignal() {
	Launch(runtime.CreateSignalShutdownChan())
}

func Launch(shutdown <-chan struct{}) {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig := config.LoadConfigOrPanic(configPrefix)
	tmanagerUrl := vConfig.GetString(sysconfig.TManagerUrlKey)
	pmanagerUrl := vConfig.GetString(sysconfig.PManagerUrlKey)

	err := runtime.CheckRequiredParams(
		fmt.Sprintf("%s.%s", configPrefix, sysconfig.TManagerUrlKey), tmanagerUrl,
		fmt.Sprintf("%s.%s", configPrefix, sysconfig.PManagerUrlKey), pmanagerUrl)
	if err != nil {
		panic(fmt.Errorf("error launching %s: %w", agentName, err))
	}

	vConfig.SetDefault(httpKey, defaultPort)
	assembler := system.NewServiceAssembler(logMonitor, vConfig, mode)

	assembler.Register(&httpclient.HttpClientServiceAssembly{})
	assembler.Register(&routing.RouterServiceAssembly{})

	assembler.Register(&client.ClientServiceAssembly{})
	assembler.Register(&job.JobServiceAssembly{})
	assembler.Register(&management.ManagementServiceAssembly{})

	runtime.AssembleAndLaunch(assembler, agentName, logMonitor, shutdown)
}
