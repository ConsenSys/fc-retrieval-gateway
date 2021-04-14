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

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientStandardCIDDiscoverRequest is used to handle client request for cid offer
func HandleClientStandardCIDDiscoverRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	pieceCID, nonce, ttl, _, _, err := fcrmessages.DecodeClientStandardDiscoverRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	now := util.GetTimeImpl().Now().Unix()
	if now > ttl {
		// Drop the connection
		return
	}

	// Search for offesr.
	offers, exists := c.OffersMgr.GetOffers(pieceCID)

	suboffers := make([]cidoffer.SubCIDOffer, 0)
	fundedPaymentChannel := make([]bool, 0)

	for _, offer := range offers {
		suboffer, err := offer.GenerateSubCIDOffer(pieceCID)
		if err != nil {
			s := "Internal error: Fail to generate suboffer."
			logging.Error(s + err.Error())
			rest.Error(w, s, http.StatusBadRequest)
			return
		}
		suboffers = append(suboffers, *suboffer)
		fundedPaymentChannel = append(fundedPaymentChannel, false) // TODO, Need to find a way to check if having payment channel set up for a given provider.
	}

	// Construct response
	response, err := fcrmessages.EncodeClientStandardDiscoverResponse(pieceCID, nonce, exists, suboffers, fundedPaymentChannel)
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
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
