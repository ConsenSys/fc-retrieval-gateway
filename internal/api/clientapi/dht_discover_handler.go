package clientapi

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
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientDHTCIDDiscoverRequest is used to handle client request for cid offer
func HandleClientDHTCIDDiscoverRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	cid, nonce, ttl, numDHT, _, _, _, err := fcrmessages.DecodeClientDHTDiscoverRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// First check if the message can be discarded
	if time.Now().Unix() > ttl {
		// Message expired.
		return
	}
	// Get a list of gatewayIDs to contact
	gateways, err := c.RegisterMgr.GetGatewaysNearCID(cid, int(numDHT))
	if err != nil || len(gateways) != int(numDHT) {
		s := "Fail to obtain required amount of peers."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	gatewayIDs := make([]*nodeid.NodeID, 0)
	for _, gateway := range gateways {
		id, err := nodeid.NewNodeIDFromHexString(gateway.NodeID)
		if err != nil {
			s := "Fail to generate node id."
			logging.Error(s + err.Error())
			rest.Error(w, s, http.StatusBadRequest)
			return
		}
		gatewayIDs = append(gatewayIDs, id)
	}

	// Construct response
	// TODO: Right now, it ignores the incremental result filed.
	// Will return all in one message.
	// Now requesting gateways.
	contacted := make([]fcrmessages.FCRMessage, 0)
	unContactable := make([]nodeid.NodeID, 0)
	for _, id := range gatewayIDs {
		res, err := c.P2PServer.RequestGatewayFromGateway(id, fcrmessages.GatewayDHTDiscoverRequestType, cid, id)
		if err != nil {
			unContactable = append(unContactable, *id)
		} else {
			contacted = append(contacted, *res)
		}
	}

	response, err := fcrmessages.EncodeClientDHTDiscoverResponse(contacted, unContactable, nonce)
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
