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
	"fmt"
	"github.com/metaform/cfm-fulcrum/internal/sysconfig"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	FulcrumClientKey system.ServiceType = "client:FulcrumClient"
	ApiClientKey     system.ServiceType = "client:ApiClient"
	fulcrumUri                          = "fulcrum.uri"
	fulcrumToken                        = "fulcrum.token"
)

type ClientServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (a *ClientServiceAssembly) Name() string {
	return "HTTP Fulcrum Client"
}

func (d *ClientServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{FulcrumClientKey, ApiClientKey}
}

func (a *ClientServiceAssembly) Init(ctx *system.InitContext) error {
	uri := ctx.Config.GetString(fulcrumUri)
	token := ctx.Config.GetString(fulcrumToken)
	tmanagerUrl := ctx.Config.GetString(sysconfig.TManagerUrlKey)
	pmanagerUrl := ctx.Config.GetString(sysconfig.PManagerUrlKey)

	err := runtime.CheckRequiredParams(
		fmt.Sprintf("%s.%s", ctx.Config.GetEnvPrefix(), fulcrumUri), uri,
		fmt.Sprintf("%s.%s", ctx.Config.GetEnvPrefix(), fulcrumToken), token,
		fmt.Sprintf("%s.%s", ctx.Config.GetEnvPrefix(), sysconfig.TManagerUrlKey), tmanagerUrl,
		fmt.Sprintf("%s.%s", ctx.Config.GetEnvPrefix(), sysconfig.PManagerUrlKey), pmanagerUrl)
	if err != nil {
		panic(fmt.Errorf("error launching %s: %w", a.Name(), err))
	}

	fulcrumClient := NewHTTPFulcrumClient(uri, token)
	ctx.Registry.Register(FulcrumClientKey, fulcrumClient)

	apiClient := NewApiClient(pmanagerUrl, fulcrumUri, "not-used")
	ctx.Registry.Register(ApiClientKey, *apiClient)

	return nil
}
