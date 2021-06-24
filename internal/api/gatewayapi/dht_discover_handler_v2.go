package gatewayapi

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
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// HandleGatewayDHTDiscoverRequestV2 handles the gateway dht discover request
func HandleGatewayDHTDiscoverRequestV2(_ *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	c := core.GetSingleInstance()

	gatewayID, pieceCID, nonce, ttl, paymentChannelAddress, voucher, err := fcrmessages.DecodeGatewayDHTDiscoverRequestV2(request)
	if err != nil {
		// Reply with invalid message
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Get the gateway's signing key
	gatewayInfo := c.RegisterMgr.GetGateway(gatewayID)
	if gatewayInfo == nil {
		logging.Warn("Gateway information not found for %s.", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	pubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		logging.Warn("Fail to obtain the public key for %s", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// First verify the message
	if request.Verify(pubKey) != nil {
		logging.Warn("Fail to verify the request from %s", gatewayID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Second check if the message can be discarded.
	if time.Now().Unix() > ttl {
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	amount, err := c.PaymentMgr.Receive(paymentChannelAddress, voucher)
	if err != nil {
		logging.Error("Payment manager receive error " + err.Error())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Respond to the request
	offers, exists := c.OffersMgr.GetOffers(pieceCID)

	subCIDOfferDigests := make([][cidoffer.CIDOfferDigestSize]byte, 0)
	fundedPaymentChannel := make([]bool, 0)

	for _, offer := range offers {
		subCIDOfferDigests = append(subCIDOfferDigests, offer.GetMessageDigest())
		fundedPaymentChannel = append(fundedPaymentChannel, false)
	}
	var response *fcrmessages.FCRMessage
	var encodingErr error
	if amount.Cmp(c.Settings.SearchPrice) < 0 {
		// not good - payment required
		logging.Error("Insufficient Funds, received " + amount.String() + ", expected: " + c.Settings.SearchPrice.String())
		// TODO get real payment channel ID
		var paymentChannelID = int64(42)
		// Construct response with payment required
		response, encodingErr = fcrmessages.EncodeGatewayDHTDiscoverResponseV2(pieceCID, nonce, exists, subCIDOfferDigests, fundedPaymentChannel, true, paymentChannelID)
	} else {
		// all good - Construct response
		response, encodingErr = fcrmessages.EncodeGatewayDHTDiscoverResponseV2(pieceCID, nonce, exists, subCIDOfferDigests, fundedPaymentChannel, false, 0)
	}
	if encodingErr != nil {
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
