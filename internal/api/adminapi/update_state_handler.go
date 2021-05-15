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
	"golang.org/x/tools/container/intsets"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// TODO: deprecated
// igor CALL fc-retrieval-gateway ->  gatewayapi -> func RequestListCIDOffer(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error)

// GatewayAdminUpdateGatewayGroupCIDOfferSupportRequest handles updating state of the Gateway, namely if it supports group CID offers
func GatewayAdminUpdateGatewayGroupCIDOfferSupportRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	c := core.GetSingleInstance()
	gateway, providers, err := fcrmessages.DecodeUpdateGatewayGroupCIDOfferSupportRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	c.GatewayID = gateway
	c.GroupCIDOfferSupportedForProviders = providers

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminInitialiseKeyResponse(true)
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}

	// Sign message
	if response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	// Send message
	w.WriteJson(response)

	go notifyProvidersOnSupportedGroupCIDOffer(c.RegisterMgr.GetAllProviders(), c.P2PServer)
}

func notifyProvidersOnSupportedGroupCIDOffer(providers []register.ProviderRegister, p2pServer *fcrp2pserver.FCRP2PServer) {
	//TODO: figure out cidMin, cidMax
	cidMin, cidMax := intsets.MinInt, intsets.MaxInt
	for _, pvd := range providers {
		id, err := nodeid.NewNodeIDFromHexString(pvd.NodeID)
		if err != nil {
			logging.Error("Error in generating node id")
			continue
		}
		go p2pServer.RequestProvider(id, fcrmessages.GatewayListDHTOfferRequestType, cidMin, cidMax, id)
	}
}
