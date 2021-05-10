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
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// RequestGatewayPing is used to request a DHT CID Discover.
func RequestGatewayPing(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	gatewayID, nonce, ttl, err := getParameters(args)
	if err != nil {
		return nil, err
	}

	// Get the core structure
	c := core.GetSingleInstance()

	// Second check if the message can be discarded.
	if time.Now().Unix() > ttl {
		return nil, err
	}

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
		return nil, err
	}
	// Get a response
	response, err := reader.Read(c.Settings.TCPInactivityTimeout)
	if err != nil {
		return nil, err
	}

	// Verify the response
	// Get the gateway's signing key
	gatewayInfo := c.RegisterMgr.GetGateway(gatewayID)
	if gatewayInfo == nil {
		return nil, errors.New("Gateway information not found")
	}

	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		return nil, errors.New("Fail to obatin the public key")
	}

	if response.Verify(pubKey) != nil {
		return nil, errors.New("Fail to verify the response")
	}

	return response, nil
}

// getParameters from args
func getParameters(args []interface{}) (*nodeid.NodeID, int64, int64, error) {
	if len(args) != 3 {
		return nil, 0, 0, errors.New("Wrong arguments")
	}
	gatewayID, ok := args[0].(*nodeid.NodeID)
	if !ok {
		return nil, 0, 0, errors.New("wrong arguments")
	}

	nonce, ok := args[1].(int64)
	if !ok {
		return nil, 0, 0, errors.New("wrong arguments")
	}

	ttl, ok := args[2].(int64)
	if !ok {
		return nil, 0, 0, errors.New("wrong arguments")
	}

	return gatewayID, nonce, ttl, nil
}
