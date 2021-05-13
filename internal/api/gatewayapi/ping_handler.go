package gatewayapi

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
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// HandleGatewayPingRequest handles admin force refresh request
func HandleGatewayPingRequest(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	logging.Debug("HandleGatewayPingRequest start ")
	defer logging.Debug("HandleGatewayPingRequest end ")
	// Get core structure
	c := core.GetSingleInstance()

	gatewayID, nonce, ttl, err := fcrmessages.DecodeGatewayPingRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		// Reply with invalid message
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Get the gateway's signing key
	gatewayInfo := c.RegisterMgr.GetGateway(gatewayID)
	if gatewayInfo == nil {
		logging.Warn("Gateway information not found for %s.", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		logging.Warn("Fail to obtain the public key for %s", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// First verify the message
	if err := request.Verify(pubKey); err != nil {
		logging.Error("Fail to verify the request from %s, error: %v", gatewayID.ToString(), err)
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Second check if the message can be discarded.
	if time.Now().Unix() > ttl {
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Respond to the request
	// if receives the request then is alive
	isAlive := true
	logging.Debug("EncodeGatewayPingResponse starting ")
	// Construct response
	response, err := fcrmessages.EncodeGatewayPingResponse(nonce, isAlive)
	if err != nil {
		logging.Debug("EncodeGatewayPingResponse error ")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	logging.Debug("EncodeGatewayPingResponse end ")

	logging.Debug("Sign %v -- %v ", c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	// Sign response
	err = response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	if err != nil {
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)

	}

	// send response
	err = writer.Write(response, c.Settings.TCPInactivityTimeout)

	return err
}
