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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleGatewayAdminInitialiseKeyRequest handles admin initilise key request
func HandleGatewayAdminInitialiseKeyRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get the core structure
	c := core.GetSingleInstance()

	nodeID, privKey, privKeyVer, err := fcrmessages.DecodeGatewayAdminInitialiseKeyRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	c.GatewayID = nodeID
	c.GatewayPrivateKey = privKey
	c.GatewayPrivateKeyVersion = privKeyVer

	// Construct messaqe
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

	// Send list cid offers
	go func() {
		// Get cid mean
		cidM, err := cid.NewContentIDFromHexString(nodeID.ToString())
		if err != nil {
			logging.Error("Error generating cid mean")
			return
		}
		gws, err := c.RegisterMgr.GetGatewaysNearCID(cidM, 16) //TODO: Do we use 16 here?
		if err != nil {
			logging.Error("Error getting gateways near cid")
			return
		}
		if len(gws) < 2 {
			logging.Error("Not enough gateways")
			return
		}
		cidMin, err := cid.NewContentIDFromHexString(gws[0].NodeID)
		if err != nil {
			logging.Error("Error generating cid min")
			return
		}
		cidMax, err := cid.NewContentIDFromHexString(gws[len(gws)-1].NodeID)
		if err != nil {
			logging.Error("Error generating cid max")
			return
		}
		pvds := c.RegisterMgr.GetAllProviders()
		for _, pvd := range pvds {
			id, err := nodeid.NewNodeIDFromHexString(pvd.NodeID)
			if err != nil {
				logging.Error("Error in generating node id")
				continue
			}
			go c.P2PServer.RequestProvider(id, fcrmessages.GatewayListDHTOfferRequestType, cidMin, cidMax, id)
		}
	}()
}
