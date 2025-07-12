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
	"os"
	"testing"
	"time"
)

const (
	testTimeout = 30 * time.Second
	streamName  = "cfm-activity"
)

func TestTestAgent_Integration(t *testing.T) {
	// Required agent config
	_ = os.Setenv("CFM-AGENT_TMANAGER_URL", "http://todo")
	_ = os.Setenv("CFM-AGENT_PMANAGER_URL", "http://todo")
	_ = os.Setenv("CFM-AGENT_FULCRUM_URI", "uri")
	_ = os.Setenv("CFM-AGENT_FULCRUM_TOKEN", "token")

	// Create and start the test agent
	shutdownChannel := make(chan struct{})
	go func() {
		Launch(shutdownChannel)
	}()

	// shut agent down
	shutdownChannel <- struct{}{}
}
