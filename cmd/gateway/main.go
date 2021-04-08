package main

// Copyright (C) 2020 ConsenSys Software Inc

import (
	_ "github.com/joho/godotenv/autoload"
	"strings"

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/config"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/adminapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

func main() {
	conf := config.NewConfig()
	appSettings := config.Map(conf)
	logging.Init(conf)
	logging.Info("Filecoin Gateway Start-up: Started")

	logging.Info("Settings: %+v", appSettings)

	// Initialise a dummy gateway instance.
	g := gateway.GetSingleInstance(&appSettings)

	// Get all registerd Gateways
	gateways, err := register.GetRegisteredGateways(appSettings.RegisterAPIURL)
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
	}
	g.RegisteredGatewaysMapLock.Lock()
	logging.Info("All registered gateways: %+v", gateways)
	for _, gateway := range gateways {
		g.RegisteredGatewaysMap[strings.ToLower(gateway.NodeID)] = &gateway
	}
	g.RegisteredGatewaysMapLock.Unlock()

	err = clientapi.StartClientRestAPI(appSettings)
	if err != nil {
		logging.Error("Error starting server: Client REST API: %s", err.Error())
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

	err = adminapi.StartAdminAPI(appSettings, g)
	if err != nil {
		logging.Error("Error starting admin tcp server: %s", err.Error())
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

	logging.Error("graceful shutdown code not written yet!")
	// TODO

	logging.Info("Filecoin Gateway Shutdown: Completed")
}
