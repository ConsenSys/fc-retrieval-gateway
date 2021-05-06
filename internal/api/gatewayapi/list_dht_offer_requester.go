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
	"errors"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// RequestListCIDOffer is used at start-up to request a list of DHT Offers from a provider with a given provider id.
func RequestListCIDOffer(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	if len(args) != 3 {
		return nil, errors.New("Wrong arguments")
	}
	cidMin, ok := args[0].(*cid.ContentID)
	if !ok {
		return nil, errors.New("Wrong arguments")
	}
	cidMax, ok := args[1].(*cid.ContentID)
	if !ok {
		return nil, errors.New("Wrong arguments")
	}
	providerID, ok := args[2].(*nodeid.NodeID)
	if !ok {
		return nil, errors.New("Wrong arguments")
	}

	// Get the core structure
	c := core.GetSingleInstance()

	request, err := fcrmessages.EncodeGatewayListDHTOfferRequest(
		c.GatewayID,
		cidMin,
		cidMax,
		c.RegistrationBlockHash,
		c.RegistrationTransactionReceipt,
		c.RegistrationMerkleRoot,
		c.RegistrationMerkleProof,
	)
	if err != nil {
		return nil, err
	}
	// Sign the request
	if request.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		return nil, errors.New("Internal error in signing the request")
	}
	// Send the request
	err = writer.Write(request, c.Settings.TCPInactivityTimeout)
	if err != nil {
		return nil, err
	}
	// Get a response
	response, err := reader.Read(c.Settings.TCPInactivityTimeout)
	if err != nil {
		return nil, err
	}

	// Verify the response
	// Get the provider's signing key
	providerInfo := c.RegisterMgr.GetProvider(providerID)
	if providerInfo == nil {
		return nil, errors.New("Provider information not found")
	}
	pubKey, err := providerInfo.GetSigningKey()
	if err != nil {
		return nil, errors.New("Fail to obatin the public key")
	}

	if response.Verify(pubKey) != nil {
		return nil, errors.New("Fail to verify the response")
	}

	// Sending acknowledgement
	// Decode the response
	cidOffers, err := fcrmessages.DecodeGatewayListDHTOfferResponse(response)
	if err != nil {
		return nil, err
	}
	// Construct the message
	cidOfferAcks := make([]fcrmessages.FCRMessage, 0)
	for _, cidOffer := range cidOffers {
		_, nonce, _, err := fcrmessages.DecodeProviderPublishDHTOfferRequest(&cidOffer)
		if err != nil {
			return nil, err
		}
		// TODO: Need to store the received offers.
		// Sign the offer message
		sig, err := fcrcrypto.SignMessage(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion, cidOffer)
		if err != nil {
			return nil, err
		}
		cidOfferAck, err := fcrmessages.EncodeProviderPublishDHTOfferResponse(nonce, sig)
		if err != nil {
			return nil, err
		}
		cidOfferAcks = append(cidOfferAcks, *cidOfferAck)
	}
	ack, err := fcrmessages.EncodeGatewayListDHTOfferAck(cidOfferAcks)
	if err != nil {
		return nil, err
	}
	// Sign the ack
	if ack.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		return nil, errors.New("Error in signing the ack")
	}

	return nil, writer.Write(ack, c.Settings.TCPInactivityTimeout)
}