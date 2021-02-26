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
	"net"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/reputation"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleAdminSetReputationChallenge(conn net.Conn, request *fcrmessages.FCRMessage) error {

	logging.Info("In handleAdminSetReputationChallenge")

	clientID, reputataion, err := fcrmessages.DecodeAdminSetReputationChallenge(request)
	if err != nil {
		return err
	}

	// Get reputation db
	rep := reputation.GetSingleInstance()
	exists := rep.ClientExists(clientID)
	var currentRep int64 = 0

	if exists {
		rep.SetClientReputation(clientID, reputataion)
		currentRep = reputataion
	}

	// Construct messaqe
	response, err := fcrmessages.EncodeAdminSetReputationResponse(clientID, currentRep, exists)
	if err != nil {
		return err
	}

	logging.Info("Admin action: Set reputation %d for client %s", currentRep, clientID)
	// Send message
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
