package gatewayapi

import (
	"math/big"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

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

// HandleGatewayDHTOfferRequest handles the gateway dht discover request
func HandleGatewayDHTOfferRequest(_ *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	c := core.GetSingleInstance()

	// TODO, Need to have an id
	pieceCID, nonce, offerDigests, paymentChannelAddress, voucher, err := fcrmessages.DecodeGatewayDHTDiscoverOfferRequest(request)
	if err != nil {
		// Reply with invalid message
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// // Get the gateway's signing key
	// gatewayInfo := c.RegisterMgr.GetGateway(gatewayID)
	// if gatewayInfo == nil {
	// 	logging.Warn("Gateway information not found for %s.", gatewayID.ToString())
	// 	return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	// }
	// pubKey, err := gatewayInfo.GetSigningKey()
	// if err != nil {
	// 	logging.Warn("Fail to obtain the public key for %s", gatewayID.ToString())
	// 	return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	// }

	// // First verify the message
	// if request.Verify(pubKey) != nil {
	// 	logging.Warn("Fail to verify the request from %s", gatewayID.ToString())
	// 	return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	// }

	amount, err := c.PaymentMgr.Receive(paymentChannelAddress, voucher)
	if err != nil {
		logging.Error("Internal error in payment manager Receive.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	lenOffers := big.NewInt(int64(len(offerDigests)))
	expectedAmount := c.Settings.OfferPrice.Mul(c.Settings.OfferPrice, lenOffers)
	if amount.Cmp(expectedAmount) < 0 {
		logging.Error("Insufficient Funds, received " + amount.String() + ", expected: " + expectedAmount.String())
		// TODO update paymentChannelID
		return writer.WriteInsufficientFunds(c.Settings.TCPInactivityTimeout, 42)
	}

	subOffers := make([]cidoffer.SubCIDOffer, len(offerDigests))
	fundedPaymentChannel := make([]bool, len(offerDigests))
	const found = true

	for i, digest := range offerDigests {
		offer, exist := c.OffersMgr.GetOfferByDigest(digest)
		fundedPaymentChannel[i] = exist

		cidOffer, err := offer.GenerateSubCIDOffer(pieceCID)
		if err != nil {
			continue
		}

		subOffers[i] = *cidOffer
	}

	// Construct response
	response, err := fcrmessages.EncodeGatewayDHTDiscoverOfferResponse(pieceCID, nonce, found, subOffers, fundedPaymentChannel)

	if err != nil {
		// TODO: Do we need a response of internal error?
		// There are three possible errors, 1. Protocol errors (request is not correct) 2. Communication errors (lost connection) and 3. Internal errors.
		// Need to do error management.
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Sign the response
	if response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		logging.Error("Internal error in signing message.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	return writer.Write(response, c.Settings.TCPInactivityTimeout)
}
