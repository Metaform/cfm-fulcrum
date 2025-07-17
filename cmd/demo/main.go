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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/metaform/cfm-fulcrum/cmd/demo/scenario"
	"os"
)

func main() {
	// Define the action flag as a string
	actionFlag := flag.String("action", "", "Action to perform: 'onboard' or 'service'")

	flag.Parse()

	switch *actionFlag {
	case "onboard":
		runOnboard()
	case "service":
		createService()
	default:
		fmt.Println("Use -action=onboard to run the onboarding process")
		fmt.Println("Use -action=service to create a service")
	}
}

func runOnboard() {
	config, err := scenario.RunOnboardCommand()
	if err != nil {
		panic(err)
	}

	// Write config file to disk
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal config to JSON: %v", err))
	}

	err = os.WriteFile("demo-config.json", configData, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write config file: %v", err))
	}

	fmt.Println("Config file written to demo-config.json")
}

func createService() {
	// Read the config file
	configData, err := os.ReadFile("demo-config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v. Please run with -action=onboard first.", err))
	}

	var config scenario.Config

	err = json.Unmarshal(configData, &config)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse config file: %v", err))
	}

	// Validate the two parameters
	if config.ProviderId == "" {
		panic("ProviderID is missing or empty in config file")
	}

	if config.AgentID == "" {
		panic("AgentID is missing or empty in config file")
	}

	if config.ServiceGroupID == "" {
		panic("ServiceGroupID is missing or empty in config file")
	}

	err = scenario.RunCreateTenantDeploymentCommand(config)
	if err != nil {
		panic(err)
	}
}
