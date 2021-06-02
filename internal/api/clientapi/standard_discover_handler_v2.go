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

// HandleClientStandardCIDDiscoverRequestV2 is used to handle client request for cid offer
func HandleClientStandardCIDDiscoverRequestV2(writer rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	pieceCID, nonce, ttl, paymentChannelAddress, voucher, err := fcrmessages.DecodeClientStandardDiscoverRequestV2(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(writer, s, http.StatusBadRequest)
		return
	}

	now := util.GetTimeImpl().Now().Unix()
	if now > ttl {
		// Drop the connection
		return
	}

	// Search for offesr.
	offers, exists := c.OffersMgr.GetOffers(pieceCID)

	var response *fcrmessages.FCRMessage

	receive, err := c.PaymentMgr.Receive(paymentChannelAddress, voucher)
	if err == nil && receive.Cmp(c.Settings.SearchPrice) >= 0 {
		// success
		subOfferDigests := make([][cidoffer.CIDOfferDigestSize]byte, 0)
		fundedPaymentChannel := make([]bool, 0)

		for _, offer := range offers {
			subOfferDigests = append(subOfferDigests, offer.GetMessageDigest())
			fundedPaymentChannel = append(fundedPaymentChannel, false)
		}

		// Construct response
		response, err = fcrmessages.EncodeClientStandardDiscoverResponseV2(pieceCID, nonce, exists, subOfferDigests, fundedPaymentChannel)
	} else {
		// Insufficient Funds Response
		if err != nil {
			logging.Error("PaymentMgr receive " + err.Error())
		} else {
			logging.Error("PaymentMgr insufficient funds received " + receive.String() + " (default: " + c.Settings.SearchPrice.String() + ")")
		}
		// TODO get real payment channel ID
		response, err = fcrmessages.EncodeInsufficientFundsResponse(42)
	}

	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(writer, s, http.StatusBadRequest)

		return
	}

	// Sign message
	err = response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	if err != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(writer, s, http.StatusInternalServerError)
		return
	}
	writer.WriteJson(response)
}
