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
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// RequestGatewayDHTDiscoverOffer is used to request a DHT
func RequestGatewayDHTDiscoverOffer(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	if len(args) != 5 {
		return nil, errors.New("wrong arguments")
	}
	cid, ok := args[0].(*cid.ContentID)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	nonce, ok := args[1].(int64)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	offerDigests, ok := args[2].([][cidoffer.CIDOfferDigestSize]byte)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	paychAddr, ok := args[3].(string)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	voucher, ok := args[4].(string)
	if !ok {
		return nil, errors.New("wrong arguments")
	}

	// Get the core structure
	c := core.GetSingleInstance()

	// // Construct message
	request, err := fcrmessages.EncodeGatewayDHTDiscoverOfferRequest(cid, nonce, offerDigests, paychAddr, voucher)
	if err != nil {
		return nil, err
	}
	// Sign the request
	if request.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		return nil, errors.New("internal error in signing the request")
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

	return response, nil
}
