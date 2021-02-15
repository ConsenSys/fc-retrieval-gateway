package register

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
)

// RegisteredNode stored information of a registered node
// Some fields maybe empty, for example, network provider AP for registered provider
type RegisteredNode struct {
	NodeID              string `json:"node_id"`
	Address             string `json:"address"`
	RootSigningKey      string `json:"root_signing_key"`
	SigningKey          string `json:"signing_key"`
	RegionCode          string `json:"region_code"`
	NetworkGatewayInfo  string `json:"network_gateway_info"`
	NetworkProviderInfo string `json:"network_provider_info"`
	NetworkClientInfo   string `json:"network_client_info"`
	NetworkAdminInfo    string `json:"network_admin_info"`
}

// GetRegisteredGateways returns registered gateways
func GetRegisteredGateways(registerURL string) ([]RegisteredNode, error) {
	url := registerURL + "/registers/gateway"
	gateways := []RegisteredNode{}
	err := request.GetJSON(url, &gateways)
	if err != nil {
		return gateways, err
	}
	return gateways, nil
}

// GetRegisteredProviders returns registered providers
func GetRegisteredProviders(registerURL string) ([]RegisteredNode, error) {
	url := registerURL + "/registers/provider"
	providers := []RegisteredNode{}
	err := request.GetJSON(url, &providers)
	if err != nil {
		return providers, err
	}
	return providers, nil
}

// GetRootSigningKey gets the root signing key
func (r *RegisteredNode) GetRootSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.RootSigningKey)
}

// GetSigningKey gets the signing key
func (r *RegisteredNode) GetSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.SigningKey)
}

// RegisterProvider to register a provider
func (r *RegisteredNode) RegisterProvider(registerURL string) error {
	url := registerURL + "/registers/provider"
	err := request.SendJSON(url, r)
	if err != nil {
		return err
	}
	return nil
}

// RegisterGateway to register a gateway
func (r *RegisteredNode) RegisterGateway(registerURL string) error {
	url := registerURL + "/registers/gateway"
	err := request.SendJSON(url, r)
	if err != nil {
		return err
	}
	return nil
}
