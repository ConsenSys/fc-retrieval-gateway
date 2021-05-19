package main

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

	_ "github.com/joho/godotenv/autoload"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrrestserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"

	"github.com/ConsenSys/fc-retrieval-gateway/config"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/adminapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
)

// Start Gateway service
func main() {
	conf := config.NewConfig()
	appSettings := config.Map(conf)
	logging.Init(conf)
	logging.Info("Filecoin Gateway Start-up: Started")

	logging.Info("Settings: %+v", appSettings)

	// Initialise a dummy gateway instance.
	c := core.GetSingleInstance(&appSettings)

	// Initialise a register manager
	c.RegisterMgr = fcrregistermgr.NewFCRRegisterMgr(appSettings.RegisterAPIURL, true, true, 10*time.Second)

	// Start register manager's routine
	c.RegisterMgr.Start()

	// Create REST Server
	c.RESTServer = fcrrestserver.NewFCRRESTServer(
		[]string{appSettings.BindAdminAPI, appSettings.BindRestAPI})

	// Add handlers to the REST Server
	c.RESTServer.
		// client api
		AddHandler(appSettings.BindRestAPI, fcrmessages.ClientEstablishmentRequestType, clientapi.HandleClientEstablishmentRequest).
		AddHandler(appSettings.BindRestAPI, fcrmessages.ClientDHTDiscoverRequestType, clientapi.HandleClientDHTCIDDiscoverRequest).
		AddHandler(appSettings.BindRestAPI, fcrmessages.ClientStandardDiscoverRequestType, clientapi.HandleClientStandardCIDDiscoverRequest).
		// admin api
		AddHandler(appSettings.BindAdminAPI, fcrmessages.GatewayAdminInitialiseKeyRequestType, adminapi.HandleGatewayAdminInitialiseKeyRequest).
		AddHandler(appSettings.BindAdminAPI, fcrmessages.GatewayAdminGetReputationRequestType, adminapi.HandleGatewayAdminGetReputationRequest).
		AddHandler(appSettings.BindAdminAPI, fcrmessages.GatewayAdminSetReputationRequestType, adminapi.HandleGatewayAdminSetReputationRequest).
		AddHandler(appSettings.BindAdminAPI, fcrmessages.GatewayAdminForceRefreshRequestType, adminapi.HandleGatewayAdminForceRefreshRequest).
		AddHandler(appSettings.BindAdminAPI, fcrmessages.GatewayAdminListDHTOfferRequestType, adminapi.HandleGatewayAdminListDHTOffersRequest).
		AddHandler(appSettings.BindAdminAPI, fcrmessages.GatewayAdminUpdateGatewayGroupCIDOfferSupportRequestType, adminapi.HandleGatewayAdminUpdateGatewayGroupCIDOfferSupportRequest)

	// Start REST Server
	err := c.RESTServer.Start()
	if err != nil {
		logging.Error("Error starting REST server: %s", err.Error())
		return
	}

	// Create P2P Server
	c.P2PServer = fcrp2pserver.NewFCRP2PServer(
		[]string{appSettings.BindGatewayAPI, appSettings.BindProviderAPI},
		c.RegisterMgr,
		appSettings.TCPInactivityTimeout)

	// Add handlers and requesters to the P2P Server
	c.P2PServer.
		// gateway api
		AddHandler(appSettings.BindGatewayAPI, fcrmessages.GatewayDHTDiscoverRequestType, gatewayapi.HandleGatewayDHTDiscoverRequest).
		AddRequester(fcrmessages.GatewayDHTDiscoverRequestType, gatewayapi.RequestGatewayDHTDiscover).
		AddRequester(fcrmessages.GatewayListDHTOfferRequestType, gatewayapi.RequestListCIDOffer).
		AddRequester(fcrmessages.GatewayNotifyProviderGroupCIDOfferSupportedRequestType, gatewayapi.NotifyProviderGroupCIDOfferSupported).
		// provider api
		AddHandler(appSettings.BindProviderAPI, fcrmessages.ProviderPublishGroupOfferRequestType, providerapi.HandleProviderPublishGroupOfferRequest).
		AddHandler(appSettings.BindProviderAPI, fcrmessages.ProviderPublishDHTOfferRequestType, providerapi.HandleProviderPublishDHTOfferRequest)

	// Start P2P Server
	err = c.P2PServer.Start()
	if err != nil {
		logging.Error("Error starting P2P server: %s", err.Error())
		return
	}

	// Configure what should be called if Control-C is hit.
	util.SetUpCtrlCExit(gracefulExit)

	logging.Info("Filecoin Gateway Start-up Complete")

	// Wait forever.
	select {}
}

// gracefulExit handles exit
func gracefulExit() {
	logging.Info("Filecoin Gateway Shutdown: Start")

	// TODO: Add shutdown process
	logging.Error("graceful shutdown code not written yet!")

	logging.Info("Filecoin Gateway Shutdown: Completed")
}
