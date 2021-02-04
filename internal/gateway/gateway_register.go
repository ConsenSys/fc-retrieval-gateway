package gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// Registration for Gateway
func Registration(settings settings.AppSettings) {
	url := settings.RegisterAPIURL + "/registers/gateway"
	providerReg := register.GatewayRegister{
		NodeID:              settings.GatewayID,
		Address:             settings.GatewayAddress,
		NetworkGatewayInfo:  settings.GatewayNetworkInfo,
		NetworkProviderInfo: settings.ProviderNetworkInfo,
		NetworkClientInfo:   settings.ClientNetworkInfo,
		NetworkAdminInfo:    settings.AdminNetworkInfo,
		RegionCode:          settings.GatewayRegionCode,
		RootSigningKey:      settings.GatewayRootSigningKey,
		SigingKey:           settings.GatewaySigningKey,
	}

	err := request.SendJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
	}
}

// GetRegisteredGateways registered Gateway list
func GetRegisteredGateways(settings settings.AppSettings) ([]register.GatewayRegister, error) {
	url := settings.RegisterAPIURL + "/registers/gateway"
	gateways := []register.GatewayRegister{}
	err := request.GetJSON(url, &gateways)
	if err != nil {
		log.Error("%+v", err)
		return gateways, err
	}
	if len(gateways) == 0 {
		log.Warn("No gateways found")
	}
	return gateways, nil
}
