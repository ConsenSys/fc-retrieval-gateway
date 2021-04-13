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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/config"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/adminapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
)

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

	err := clientapi.StartClientRestAPI(appSettings)
	if err != nil {
		logging.Error("Error starting client REST server: %s", err.Error())
		return
	}

	err = adminapi.StartAdminAPI(appSettings)
	if err != nil {
		logging.Error("Error starting admin REST server: %s", err.Error())
		return
	}

	err = gatewayapi.StartGatewayAPI(appSettings)
	if err != nil {
		logging.Error("Error starting gateway tcp server: %s", err.Error())
		return
	}

	err = providerapi.StartProviderAPI(appSettings)
	if err != nil {
		logging.Error("Error starting provider tcp server: %s", err.Error())
		return
	}

	// Configure what should be called if Control-C is hit.
	util.SetUpCtrlCExit(gracefulExit)

	logging.Info("Filecoin Gateway Start-up Complete")

	// Wait forever.
	select {}
}

func gracefulExit() {
	logging.Info("Filecoin Gateway Shutdown: Start")

	// TODO: Add shutdown process
	logging.Error("graceful shutdown code not written yet!")

	logging.Info("Filecoin Gateway Shutdown: Completed")
}
