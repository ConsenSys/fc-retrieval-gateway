package main

// Copyright (C) 2020 ConsenSys Software Inc

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/ConsenSys/fc-retrieval-gateway/config"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/adminapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/register"
)

func main() {
	conf := config.NewConfig()
	settings := config.Map(conf)
	logging.Init(conf)
	logging.Info("Filecoin Gateway Start-up: Started")

	logging.Info("Settings: %+v", settings)

	// Initialise a dummy gateway instance.
	g := gateway.GetSingleInstance(&settings)

	// Register Gateway
	gatewayReg := register.GatewayRegister{
		NodeID:              settings.GatewayID,
		Address:             settings.GatewayAddress,
		RootSigningKey:      settings.GatewayRootSigningKey,
		SigningKey:          settings.GatewaySigningKey,
		NetworkGatewayInfo:  settings.GatewayNetworkInfo,
		NetworkProviderInfo: settings.ProviderNetworkInfo,
		NetworkClientInfo:   settings.ClientNetworkInfo,
		NetworkAdminInfo:    settings.AdminNetworkInfo,
		RegionCode:          settings.GatewayRegionCode,
	}
	gatewayReg.RegisterGateway(settings.RegisterAPIURL)

	// Get all registerd Gateways
	gateways, err := register.GetRegisteredGateways(settings.RegisterAPIURL)
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
	}
	g.RegisteredGatewaysMapLock.Lock()
	logging.Info("All registered gateways: %+v", gateways)
	for _, gateway := range gateways {
		g.RegisteredGatewaysMap[gateway.NodeID] = &gateway
	}
	g.RegisteredGatewaysMapLock.Unlock()

	// Get all registered Providers
	providers, err := register.GetRegisteredProviders(settings.RegisterAPIURL)
	if err != nil {
		logging.Error("Unable to get registered providers: %v", err)
	}
	g.RegisteredProvidersMapLock.Lock()
	logging.Info("All registered providers: %+v", providers)
	for _, provider := range providers {
		g.RegisteredGatewaysMap[provider.NodeID] = &provider
	}
	g.RegisteredGatewaysMapLock.Unlock()

	err = clientapi.StartClientRestAPI(settings)
	if err != nil {
		logging.Error("Error starting server: Client REST API: %s", err.Error())
		return
	}

	err = gatewayapi.StartGatewayAPI(settings)
	if err != nil {
		logging.Error("Error starting gateway tcp server: %s", err.Error())
		return
	}

	err = providerapi.StartProviderAPI(settings)
	if err != nil {
		logging.Error("Error starting provider tcp server: %s", err.Error())
		return
	}

	err = adminapi.StartAdminAPI(settings, g)
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
