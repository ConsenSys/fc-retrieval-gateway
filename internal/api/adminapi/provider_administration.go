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
	"strings"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleAdminEnrollProvider(conn net.Conn, request *fcrmessages.FCRMessage, settings settings.AppSettings) error {
	gatewayCoreStructure := gateway.GetSingleInstance()
	if gatewayCoreStructure.GatewayPrivateKey == nil {
		return errors.New("this gateway hasn't been initialised by the admin")
	}

	nodeID, address, rootSigningKey, signingKey, regionCode, networkInfoGateway,
		networkInfoClient, networkInfoAdmin, err :=
		fcrmessages.DecodeGatewayAdminEnrollProviderRequest(request)
	if err != nil {
		return err
	}

	newProvider := register.ProviderRegister{
		NodeID:             nodeID.ToString(),
		Address:            address,
		RootSigningKey:     rootSigningKey,
		SigningKey:         signingKey,
		RegionCode:         regionCode,
		NetworkInfoGateway: networkInfoGateway,
		NetworkInfoClient:  networkInfoClient,
		NetworkInfoAdmin:   networkInfoAdmin,
	}

	newProviderRegistered := registerProvider(newProvider, gatewayCoreStructure)
	if !newProviderRegistered {
		logging.Error("can not register a provider")
	}

	// Construct a message
	response, err := fcrmessages.EncodeGatewayAdminEnrollProviderResponse(newProviderRegistered)
	if err != nil {
		return err
	}
	// Sign message
	if response.Sign(gatewayCoreStructure.GatewayPrivateKey, gatewayCoreStructure.GatewayPrivateKeyVersion) != nil {
		return errors.New("error in signing message")
	}
	// Send message
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.TCPInactivityTimeout)
}

func registerProvider(newProvider register.ProviderRegister, gatewayCoreStructure *gateway.Gateway) (ok bool) {
	logging.Debug("Registered a new provider")
	update := false
	gatewayCoreStructure.RegisteredProvidersMapLock.RLock()
	storedInfo, exist := gatewayCoreStructure.RegisteredProvidersMap[strings.ToLower(newProvider.NodeID)]
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
		if newProvider.Address != storedInfo.GetAddress() ||
			newProvider.NetworkInfoAdmin != storedInfo.GetNetworkInfoAdmin() ||
			newProvider.NetworkInfoClient != storedInfo.GetNetworkInfoClient() ||
			newProvider.NetworkInfoGateway != storedInfo.GetNetworkInfoGateway() ||
			newProvider.RegionCode != storedInfo.GetRegionCode() ||
			newProvider.RootSigningKey != rootSigningKey ||
			newProvider.SigningKey != signingKey {

			update = true
		}
	}
	gatewayCoreStructure.RegisteredProvidersMapLock.RUnlock()
	if update {
		logging.Info("Add to registered providers map: nodeID=%+v", newProvider.NodeID)
		gatewayCoreStructure.RegisteredProvidersMapLock.Lock()
		gatewayCoreStructure.RegisteredProvidersMap[strings.ToLower(newProvider.NodeID)] = &newProvider
		gatewayCoreStructure.RegisteredProvidersMapLock.Unlock()
	}
	return true
}
