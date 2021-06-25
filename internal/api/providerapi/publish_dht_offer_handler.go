/*
Package providerapi - set of remote API used to call a Gateway, grouped to a specific caller type - Retrieval Provider.
All calls from FileCoin Secondary Retrieval network nodes of type Retrieval Provider are going to API handlers in this package.
*/
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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// HandleProviderPublishDHTOfferRequest handles the provider publish dht offer request
func HandleProviderPublishDHTOfferRequest(_ *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	c := core.GetSingleInstance()

	providerID, nonce, offers, err := fcrmessages.DecodeProviderPublishDHTOfferRequest(request)
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

	// Verify the offer one by one
	for _, offer := range offers {
		if offer.Verify(pubKey) != nil {
			logging.Warn("Fail to verify the offer from %s", providerID.ToString())
			return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
		}

		if c.OffersMgr.AddDHTOffer(&offer) != nil {
			logging.Error("Internal error in adding single cid offer.")
			return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
		}
	}

	// Sign the request
	sig, err := fcrcrypto.SignMessage(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion, request.GetMessageBody())
	if err != nil {
		logging.Error("Internal error in signing message.")
		return writer.WriteInvalidMessage(c.Settings.TCPInactivityTimeout)
	}

	// Construct response
	response, err := fcrmessages.EncodeProviderPublishDHTOfferResponse(nonce, sig)
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
