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
	"sync"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleAdminAcceptKeysChallenge(conn net.Conn, request *fcrmessages.FCRMessage, wg *sync.WaitGroup) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	nodeID, privKey, privKeyVer, err := fcrmessages.DecodeGatewayAdminInitialiseKeyRequest(request)
	if err != nil {
		return err
	}

	g.GatewayID = nodeID
	g.GatewayPrivateKey = privKey
	g.GatewayPrivateKeyVersion = privKeyVer
	wg.Done() // need mutex to protect g

	// Construct messaqe
	response, err := fcrmessages.EncodeGatewayAdminInitialiseKeyResponse(true)
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
