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
	"encoding/json"
	"net"
	"strconv"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/reputation"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleAdminGetReputationChallenge(conn net.Conn, request *messages.AdminGetReputationChallenge) error {

	logging.Info("In handleAdminGetReputationChallenge")

	// Get gateway core struct
	gw := gateway.GetSingleInstance()
	clientID := request.ClientID

	// Construct response
	response := messages.AdminGetReputationResponse{
		MessageType:     messages.AdminGetReputationResponseType,
		ProtocolVersion: gw.ProtocolVersion,
		ClientID:        request.ClientID}

	rep := reputation.GetSingleInstance()
	id, err := nodeid.NewNodeIDFromString(clientID)
	if err != nil {
		logging.Info("Cannot find clientID: %s", err)
		// Message is invalid.
		tcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout, "Message is invalid.")
	}

	response.Reputation, response.Exists = rep.GetClientReputation(id)
	// Send message
	data, _ := json.Marshal(response)
	logging.Info("Admin action: Returned reputation %s for client %s", strconv.FormatInt(response.Reputation, 10), clientID)
	return tcpcomms.SendTCPMessage(conn, messages.AdminGetReputationResponseType, data, settings.DefaultTCPInactivityTimeout)
}
