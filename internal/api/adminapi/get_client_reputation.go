package adminapi

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"errors"
	"net"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/reputation"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleAdminGetReputationChallenge(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get core structure
	g := gateway.GetSingleInstance()
	if g.GatewayPrivateKey == nil {
		return errors.New("This gateway hasn't been initialised by the admin")
	}

	clientID, err := fcrmessages.DecodeGatewayAdminGetReputationRequest(request)
	if err != nil {
		return err
	}

	// Get reputation db
	rep := reputation.GetSingleInstance()
	reputation, exists := rep.GetClientReputation(clientID)

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminGetReputationResponse(clientID, reputation, exists)
	if err != nil {
		return err
	}
	// Sign message
	if response.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		return errors.New("Error in signing message")
	}
	// Send message
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
