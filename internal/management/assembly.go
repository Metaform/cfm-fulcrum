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

package management

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/metaform/cfm-fulcrum/internal/client"
	"github.com/metaform/connector-fabric-manager/assembly/httpclient"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/system"
	"net/http"
)

type ManagementServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (a *ManagementServiceAssembly) Name() string {
	return "Management API"
}

func (d *ManagementServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{routing.RouterKey, httpclient.HttpClientKey}
}

func (a *ManagementServiceAssembly) Init(context *system.InitContext) error {
	router := context.Registry.Resolve(routing.RouterKey).(chi.Router)
	fulcrumClient := context.Registry.Resolve(client.FulcrumClientKey).(client.FulcrumClient)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		response := response{Message: "OK"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	router.Post("/fulcrum-token", func(w http.ResponseWriter, r *http.Request) {
		var result map[string]any
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, fmt.Sprintf("failed to unmarshal JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Extract token from the request body
		token, ok := result["token"].(string)
		if !ok {
			http.Error(w, "token field is required and must be a string", http.StatusBadRequest)
			return
		}

		if err := fulcrumClient.UpdateToken(token); err != nil {
			context.LogMonitor.Severef("error updating token: %w", err)
			http.Error(w, fmt.Sprintf("error updating token: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		response := response{Message: "OK"}
		json.NewEncoder(w).Encode(response)
	})

	return nil

}

type response struct {
	Message string `json:"message"`
}
