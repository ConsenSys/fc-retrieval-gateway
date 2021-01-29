package gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
)

// Register data model
type Register struct {
	Address        string
	NetworkInfo    string
	RegionCode     string
	RootSigningKey string
	SigingKey      string
}

// Registration registers a gateway
func Registration(url string, settings settings.AppSettings) {

	providerReg := Register{
		Address:        settings.GatewayAddress,
		NetworkInfo:    settings.GatewayNetworkInfo,
		RegionCode:     settings.GatewayRegionCode,
		RootSigningKey: settings.GatewayRootSigningKey,
		SigingKey:      settings.GatewaySigningKey,
	}

	err := request.SendJSON(url, providerReg)
	if err != nil {
		log.Error("%+v", err)
	}
}
