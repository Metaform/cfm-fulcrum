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
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	ClientKey    system.ServiceType = "client:Client"
	fulcrumUri                      = "fulcrum.uri"
	fulcrumToken                    = "fulcrum.token"
)

type ClientServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (a *ClientServiceAssembly) Name() string {
	return "HTTP Fulcrum Client"
}

func (d *ClientServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{ClientKey}
}

func (a *ClientServiceAssembly) Init(context *system.InitContext) error {
	uri := context.Config.GetString(fulcrumUri)
	token := context.Config.GetString(fulcrumToken)
	err := runtime.CheckRequiredParams(
		fmt.Sprintf("%s.%s", context.Config.GetEnvPrefix(), fulcrumUri), uri,
		fmt.Sprintf("%s.%s", context.Config.GetEnvPrefix(), fulcrumToken), token)
	if err != nil {
		panic(fmt.Errorf("error launching %s: %w", a.Name(), err))
	}

	client := NewHTTPFulcrumClient(uri, token)
	context.Registry.Register(ClientKey, client)
	return nil
}
