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

	// Check amount
	amount, err := c.PaymentMgr.Receive(paymentChannel, voucher)
	if err != nil {
		s := "Internal error in payment manager Receive."
		logging.Error(s)
		logging.Error(err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	unit := 0
	for _, entry := range allGatewaysOfferDigests {
		unit += len(entry)
	}
	expectedAmt := new(big.Int).Mul(big.NewInt(int64(unit)), c.Settings.OfferPrice)
	if amount.Cmp(expectedAmt) < 0 {
		s := "Insufficient Funds, received " + amount.String() + ", expected: " + expectedAmt.String()
		logging.Error(s)
		rest.Error(w, s, http.StatusInternalServerError)
	}

	contactedGateways := make([]nodeid.NodeID, 0)
	contactedResp := make([]fcrmessages.FCRMessage, 0)
	unContactable := make([]nodeid.NodeID, 0)
	for idx, targetGatewayID := range targetGatewayIDs {
		thisGatewayOfferDigests := allGatewaysOfferDigests[idx]
		targetGateway := c.RegisterMgr.GetGateway(&targetGatewayID)
		// Pay this gateway
		toPay := new(big.Int).Mul(big.NewInt(int64(len(thisGatewayOfferDigests))), c.Settings.OfferPrice)
		paychAddr, voucher, topup, err := c.PaymentMgr.Pay(targetGateway.GetAddress(), 0, toPay)
		if err != nil {
			s := "Fail to pay recipient."
			logging.Error(s + err.Error())
			rest.Error(w, s, http.StatusBadRequest)
			return
		}
		if topup {
			if err := c.PaymentMgr.Topup(targetGateway.GetAddress(), c.Settings.TopupAmount); err != nil {
				s := "Fail to top up."
				logging.Error(s + err.Error())
				rest.Error(w, s, http.StatusBadRequest)
				return
			}
			paychAddr, voucher, topup, err = c.PaymentMgr.Pay(targetGateway.GetAddress(), 0, c.Settings.SearchPrice)
			if err != nil {
				s := "Fail to pay recipient."
				logging.Error(s + err.Error())
				rest.Error(w, s, http.StatusBadRequest)
				return
			}
		}
		// using index from one collection to access another; create a struct?
		res, err := c.P2PServer.RequestGatewayFromGateway(&targetGatewayID, fcrmessages.GatewayDHTDiscoverOfferRequestType, cid, &targetGatewayID, nonce, thisGatewayOfferDigests, paychAddr, voucher)
		if err != nil {
			logging.Info("Uncontactable: %v", err.Error())
			unContactable = append(unContactable, targetGatewayID)
		} else {
			contactedGateways = append(contactedGateways, targetGatewayID)
			contactedResp = append(contactedResp, *res)
		}
	}

	response, err := fcrmessages.EncodeClientDHTDiscoverOfferResponse(cid, nonce, contactedGateways, contactedResp, false, 0)
	if err != nil {
		s := "Internal error: Fail to encode message, type: " + strconv.Itoa(fcrmessages.ClientDHTDiscoverOfferResponseType)
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}

	// Sign message
	if signErr := response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion); signErr != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + signErr.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}

	if writeErr := w.WriteJson(response); writeErr != nil {
		logging.Error("can't write JSON during HandleClientDHTDiscoverOfferRequest %s", writeErr.Error())
	}
}
