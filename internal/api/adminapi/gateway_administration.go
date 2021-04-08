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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
	"net"
	"strings"
)

func handleAdminEnrollGateway(conn net.Conn, request *fcrmessages.FCRMessage, settings settings.AppSettings) error {
	// Get core structure
	g := gateway.GetSingleInstance()
	if g.GatewayPrivateKey == nil {
		return errors.New("this gateway hasn't been initialised by the admin")
	}

	nodeID, address, rootSigningKey, signingKey, regionCode, networkInfoGateway, networkInfoProvider,
		networkInfoClient, networkInfoAdmin, err :=
		fcrmessages.DecodeGatewayAdminEnrollGatewayRequest(request)
	if err != nil {
		return err
	}

	newGateway := register.GatewayRegister{
		NodeID:              nodeID.ToString(),
		Address:             address,
		RootSigningKey:      rootSigningKey,
		SigningKey:          signingKey,
		RegionCode:          regionCode,
		NetworkInfoGateway:  networkInfoGateway,
		NetworkInfoProvider: networkInfoProvider,
		NetworkInfoClient:   networkInfoClient,
		NetworkInfoAdmin:    networkInfoAdmin,
	}
	ok := registerGateway(newGateway)
	if !ok {
		return errors.New("can not register a gateway")
	}

	// Construct a message
	response, err := fcrmessages.EncodeGatewayAdminEnrollGatewayRequest(
		nodeID,
		address,
		rootSigningKey,
		signingKey,
		regionCode,
		networkInfoGateway,
		networkInfoProvider,
		networkInfoClient,
		networkInfoAdmin)
	if err != nil {
		return err
	}
	// Sign message
	if response.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		return errors.New("error in signing message")
	}
	// Send message
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.TCPInactivityTimeout)
}

func registerGateway(newGateway register.GatewayRegister) (ok bool) {
	logging.Debug("Register a new gateway")

	// Get core structure
	g := gateway.GetSingleInstance()
	if g.GatewayPrivateKey == nil {
		logging.Error("This gateway hasn't been initialised by the admin")
		return false
	}

	update := false
	g.RegisteredGatewaysMapLock.RLock()
	storedInfo, exist := g.RegisteredGatewaysMap[strings.ToLower(newGateway.NodeID)]
	if !exist {
		update = true
	} else {
		key, err := storedInfo.GetRootSigningKey()
		rootSigningKey, err2 := key.EncodePublicKey()
		key, err3 := storedInfo.GetSigningKey()
		signingKey, err4 := key.EncodePublicKey()
		if err != nil || err2 != nil || err3 != nil || err4 != nil {
			logging.Error("Error in generating key string")
			return false
		}
		if newGateway.Address != storedInfo.GetAddress() ||
			newGateway.NetworkInfoAdmin != storedInfo.GetNetworkInfoAdmin() ||
			newGateway.NetworkInfoClient != storedInfo.GetNetworkInfoClient() ||
			newGateway.NetworkInfoProvider != storedInfo.GetNetworkInfoProvider() ||
			newGateway.NetworkInfoGateway != storedInfo.GetNetworkInfoGateway() ||
			newGateway.RegionCode != storedInfo.GetRegionCode() ||
			newGateway.RootSigningKey != rootSigningKey ||
			newGateway.SigningKey != signingKey {

			update = true
		}
	}
	g.RegisteredGatewaysMapLock.RUnlock()
	if update {
		logging.Info("Add to registered gateways map: nodeID=%+v", newGateway.NodeID)
		g.RegisteredGatewaysMapLock.Lock()
		g.RegisteredGatewaysMap[strings.ToLower(newGateway.NodeID)] = &newGateway
		g.RegisteredGatewaysMapLock.Unlock()
	}
	return true
}
