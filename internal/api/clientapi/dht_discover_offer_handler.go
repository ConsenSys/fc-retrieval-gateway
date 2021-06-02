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
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

func HandleClientDHTDiscoverOfferRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	cid, nonce, allGatewaysOfferDigests, targetGatewayIDs, paymentChannel, voucher, err := fcrmessages.DecodeClientDHTDiscoverOfferRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	contactedGateways := make([]nodeid.NodeID, 0)
	contactedResp := make([]fcrmessages.FCRMessage, 0)
	unContactable := make([]nodeid.NodeID, 0)
	for idx, targetGatewayID := range targetGatewayIDs {
		// using index from one collection to access another; create a struct?
		thisGatewayOfferDigests := allGatewaysOfferDigests[idx]
		res, err := c.P2PServer.RequestGatewayFromGateway(&targetGatewayID, fcrmessages.GatewayDHTDiscoverOfferRequestType, cid, nonce, thisGatewayOfferDigests, paymentChannel, voucher)
		if err != nil {
			unContactable = append(unContactable, targetGatewayID)
		} else {
			contactedGateways = append(contactedGateways, targetGatewayID)
			contactedResp = append(contactedResp, *res)
		}
	}

	response, err := fcrmessages.EncodeClientDHTDiscoverOfferResponse(cid, nonce, contactedGateways, contactedResp)
	if err != nil {
		s := "Internal error: Fail to encode message, type: " + strconv.Itoa(fcrmessages.ClientDHTDiscoverOfferResponseType)
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
