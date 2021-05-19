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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

func NotifyProviderGroupCIDOfferSupported(reader *fcrp2pserver.FCRServerReader, writer *fcrp2pserver.FCRServerWriter, args ...interface{}) (*fcrmessages.FCRMessage, error) {
	// Get parameters
	if len(args) != 1 {
		return nil, errors.New("wrong arguments")
	}
	groupCIDOfferSupported, ok := args[0].(bool)
	if !ok {
		return nil, errors.New("wrong arguments")
	}

	// Get the core structure
	c := core.GetSingleInstance()

	request, err := fcrmessages.EncodeGatewayNotifyProviderGroupCIDOfferSupportRequest(
		c.GatewayID,
		groupCIDOfferSupported,
	)
	if err != nil {
		return nil, err
	}

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

	acknowledged, err := fcrmessages.DecodeGatewayNotifyProviderGroupCIDOfferSupportResponse(response)
	if err != nil {
		return nil, err
	}
	if !acknowledged {
		return nil, errors.New("NotifyProviderGroupCIDOfferSupported: no successful acknowledgement from Provider to Gateway")
	}
	return response, nil
}
