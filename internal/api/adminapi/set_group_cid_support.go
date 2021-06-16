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
  "net/http"

  "github.com/ant0ine/go-json-rest/rest"

  "github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
  "github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
  "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
  "github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
  "github.com/ConsenSys/fc-retrieval-common/pkg/register"
  "github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// HandleGatewayAdminUpdateGatewayGroupCIDOfferSupportRequest handles updating state of the Gateway, namely if it supports group CID offers
func HandleGatewayAdminUpdateGatewayGroupCIDOfferSupportRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	c := core.GetSingleInstance()
	_, providerIDs, err := fcrmessages.DecodeUpdateGatewayGroupCIDOfferSupportRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	c.GroupCIDOfferSupportedForProviders = providerIDs

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminInitialiseKeyResponse(true)
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}

	// Sign message
	err = response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	if err != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	// Send message
  if err := w.WriteJson(response); err != nil {
    logging.Error("can't write JSON during HandleGatewayAdminUpdateGatewayGroupCIDOfferSupportRequest %s", err.Error())
  }

	notifyProvidersOnSupportedGroupCIDOffer(c.RegisterMgr.GetAllProviders(), c.P2PServer, c.GatewayID)
}

func notifyProvidersOnSupportedGroupCIDOffer(providers []register.ProviderRegistrar, p2pServer *fcrp2pserver.FCRP2PServer, callerGatewayId *nodeid.NodeID) {
	for _, pvd := range providers {
		providerID, err := nodeid.NewNodeIDFromHexString(pvd.GetNodeID())
		if err != nil {
			logging.Error("Error in generating node id")
			continue
		}
		const thisGatewaySupportsGroupCIDOffer = true
    go notifyProvider(p2pServer, providerID, fcrmessages.GatewayListDHTOfferRequestType, callerGatewayId, thisGatewaySupportsGroupCIDOffer)
	}
}

func notifyProvider(p2pServer *fcrp2pserver.FCRP2PServer, providerID *nodeid.NodeID, msgType int32, gatewayID *nodeid.NodeID, thisGatewaySupportsGroupCIDOffer bool) {
  providerResponse, err := p2pServer.RequestProvider(providerID, msgType, gatewayID, thisGatewaySupportsGroupCIDOffer)
  if err != nil {
    logging.Error("error notifying provider, method: %d; provider id: %s", msgType, providerID)
  }
  logging.Info("Provider response: %v", providerResponse)
}
