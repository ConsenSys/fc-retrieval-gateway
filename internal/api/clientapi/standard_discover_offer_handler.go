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
	"math/big"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientStandardDiscoverOfferRequest is used to receive payment to respond to client standard offer query
func HandleClientStandardDiscoverOfferRequest(writer rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	pieceCID, nonce, ttl, offerDigests, paymentChannelAddress, voucher, err := fcrmessages.DecodeClientStandardDiscoverOfferRequest(request)
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

	var response *fcrmessages.FCRMessage

	receive, err := c.PaymentMgr.Receive(paymentChannelAddress, voucher)
	totalPrice := big.NewInt(int64(len(offerDigests))).Mul(big.NewInt(1), c.Settings.OfferPrice)
	if err == nil && receive.Cmp(totalPrice) >= 0 {
		// Success - Search for offers
		subOffers := make([]cidoffer.SubCIDOffer, len(offerDigests))
		fundedPaymentChannel := make([]bool, len(offerDigests))
		found := false

		for i, digest := range offerDigests {
			offer, exist := c.OffersMgr.GetOfferByDigest(digest)
			fundedPaymentChannel[i] = exist
			found = exist

			cidOffer, err := offer.GenerateSubCIDOffer(pieceCID)
			if err != nil {
				continue
			}

			subOffers[i] = *cidOffer
		}

		// Construct response
		response, err = fcrmessages.EncodeClientStandardDiscoverOfferResponse(pieceCID, nonce, found, subOffers, fundedPaymentChannel)
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
