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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleGatewayAdminListDHTOffersRequest handles admin list dht offer request
func HandleGatewayAdminListDHTOffersRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	if c.GatewayPrivateKey == nil {
		s := "This gateway hasn't been initialised by the admin"
		logging.Error(s)
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	refresh, err := fcrmessages.DecodeGatewayAdminListDHTOfferRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	refreshed := false
	if refresh {
		// Send list cid offers
		go func() {
			cidMin, cidMax, err := c.RegisterMgr.GetGatewayCIDRange(c.GatewayID)
			if err != nil {
				logging.Error("Error getting cid range for gateway initial list dht offer: %v", err.Error())
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
		refreshed = true
	}

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminListDHTOfferResponse(refreshed)
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
	w.WriteJson(response)
}
