package providerapi

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
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

// HandleProviderPublishGroupOfferRequest handles the provider publish group offer request
func HandleProviderPublishGroupOfferRequest(_ *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	c := core.GetSingleInstance()

	// TODO Add nonce, it looks like nonce is not needed
	providerID, _, offer, err := fcrmessages.DecodeProviderPublishGroupOfferRequest(request)
	if err != nil {
		// Reply with invalid message
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Verify the request
	// Get the provider's signing key
	providerInfo := c.RegisterMgr.GetProvider(providerID)
	if providerInfo == nil {
		logging.Warn("Provider information not found for %s.", providerID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	pubKey, err := providerInfo.GetSigningKey()
	if err != nil {
		logging.Warn("Fail to obtain the public key for %s", providerID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}
	if request.Verify(pubKey) != nil {
		logging.Warn("Fail to verify the request from %s", providerID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Verify the offer
	if offer.Verify(pubKey) != nil {
		logging.Warn("Fail to verify the offer from %s", providerID.ToString())
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Store the offer
	if c.OffersMgr.AddGroupOffer(offer) != nil {
		logging.Error("Internal error in adding group cid offer.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Construct the response
	response, err := fcrmessages.EncodeProviderPublishGroupOfferResponse(
		*c.GatewayID,
		offer.GetMessageDigest(),
	)
	if err != nil {
		logging.Error("Internal error in encoding message.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Sign the response
	if response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		logging.Error("Internal error in signing message.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	return writer.Write(response, c.Settings.TCPInactivityTimeout)
}
