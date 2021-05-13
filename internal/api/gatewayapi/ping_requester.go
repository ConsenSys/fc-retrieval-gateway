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
	"errors"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// RequestGatewayPing is used to request a DHT CID Discover.
func RequestGatewayPing(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	logging.Debug("RequestGatewayPing start")
	defer logging.Debug("RequestGatewayPing end")
	// Get parameters
	if len(args) != 1 {
		return nil, errors.New("wrong arguments")
	}

	gatewayID, ok := args[0].(*nodeid.NodeID)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	// FIXME
	nonce := int64(0)
	ttl := int64(0)

	// Get the core structure
	c := core.GetSingleInstance()

	// Construct message
	request, err := fcrmessages.EncodeGatewayPingRequest(gatewayID, nonce, ttl)
	if err != nil {
		return nil, err
	}

	// Sign the request
	if request.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		return nil, errors.New("Internal error in signing the request")
	}
	// Send the request
	err = writer.Write(request, c.Settings.TCPInactivityTimeout)
	if err != nil {
		logging.Error("writer.Write error %v", err)
		return nil, err
	}

	// Get a response
	response, err := reader.Read(c.Settings.TCPInactivityTimeout)
	if err != nil {
		return nil, err
	}
	logging.Debug("response  %v %v", response, err)

	// Verify the response
	// Get the gateway's signing key
	gatewayInfo := c.RegisterMgr.GetGateway(gatewayID)
	if gatewayInfo == nil {
		return nil, errors.New("Gateway information not found")
	}

	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		return nil, errors.New("fail to obtain the public key")
	}

	if err := response.Verify(pubKey); err != nil {
		logging.Debug("fail to verify the response with pubkey, error: %v", err)
		// TODO uncomment next line after sign message
		// return nil, errors.New("Fail to verify the response")
	}

	return response, nil
}
