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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

// StartGatewayAPI starts the TCP API as a separate go routine.
func StartGatewayAPI(settings settings.AppSettings) error {
	// Get core structure
	c := core.GetSingleInstance()

	// Initialise a new P2P Server
	c.GatewayServer = fcrp2pserver.NewFCRP2PServer("gateway-server", c.RegisterMgr, settings.TCPInactivityTimeout)

	// Add handlers
	c.GatewayServer.
		AddHandler(fcrmessages.GatewayDHTDiscoverRequestType, handleGatewayDHTDiscoverRequest).
		AddRequester(fcrmessages.GatewayDHTDiscoverRequestType, requestGatewayDHTDiscover).
		AddRequester(fcrmessages.GatewayListDHTOfferRequestType, requestListCIDOffers)

	// Start server
	return c.GatewayServer.Start(settings.BindGatewayAPI)
}
