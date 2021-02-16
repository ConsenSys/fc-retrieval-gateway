package register

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
)

// RegisteredNode stored network information of a registered node
type RegisteredNode interface {
	GetNodeID() string
	GetAddress() string
	GetRegionCode() string
	GetRootSigningKey() (*fcrcrypto.KeyPair, error)
	GetSigningKey() (*fcrcrypto.KeyPair, error)
	GetNetworkGatewayInfo() string
	GetNetworkProviderInfo() string
	GetNetworkClientInfo() string
	GetNetworkAdminInfo() string
}

// GatewayRegister stores information of a registered gateway
type GatewayRegister struct {
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

// ProviderRegister stores information of a registered provider
type ProviderRegister struct {
	NodeID             string `json:"node_id"`
	Address            string `json:"address"`
	RootSigningKey     string `json:"root_signing_key"`
	SigningKey         string `json:"signing_key"`
	RegionCode         string `json:"region_code"`
	NetworkGatewayInfo string `json:"network_gateway_info"`
	NetworkClientInfo  string `json:"network_client_info"`
	NetworkAdminInfo   string `json:"network_admin_info"`
}

// GetRegisteredGateways returns registered gateways
func GetRegisteredGateways(registerURL string) ([]GatewayRegister, error) {
	url := registerURL + "/registers/gateway"
	gateways := []GatewayRegister{}
	err := request.GetJSON(url, &gateways)
	if err != nil {
		return gateways, err
	}
	return gateways, nil
}

// GetRegisteredProviders returns registered providers
func GetRegisteredProviders(registerURL string) ([]ProviderRegister, error) {
	url := registerURL + "/registers/provider"
	providers := []ProviderRegister{}
	err := request.GetJSON(url, &providers)
	if err != nil {
		return providers, err
	}
	return providers, nil
}

// GetNodeID gets the node id
func (r *GatewayRegister) GetNodeID() string {
	return r.NodeID
}

// GetNodeID gets the node id
func (r *ProviderRegister) GetNodeID() string {
	return r.NodeID
}

// GetAddress gets the address
func (r *GatewayRegister) GetAddress() string {
	return r.Address
}

// GetAddress gets the node id
func (r *ProviderRegister) GetAddress() string {
	return r.Address
}

// GetRegionCode gets the region code
func (r *GatewayRegister) GetRegionCode() string {
	return r.RegionCode
}

// GetRegionCode gets the region code
func (r *ProviderRegister) GetRegionCode() string {
	return r.RegionCode
}

// GetNetworkGatewayInfo gets the network gateway ap
func (r *GatewayRegister) GetNetworkGatewayInfo() string {
	return r.NetworkGatewayInfo
}

// GetNetworkGatewayInfo gets the network gateway ap
func (r *ProviderRegister) GetNetworkGatewayInfo() string {
	return r.NetworkGatewayInfo
}

// GetNetworkProviderInfo gets the network provider ap
func (r *GatewayRegister) GetNetworkProviderInfo() string {
	return r.NetworkProviderInfo
}

// GetNetworkProviderInfo gets the network provider ap
func (r *ProviderRegister) GetNetworkProviderInfo() string {
	return ""
}

// GetNetworkClientInfo gets the network client ap
func (r *GatewayRegister) GetNetworkClientInfo() string {
	return r.NetworkClientInfo
}

// GetNetworkClientInfo gets the network client ap
func (r *ProviderRegister) GetNetworkClientInfo() string {
	return r.NetworkClientInfo
}

// GetNetworkAdminInfo gets the network admin ap
func (r *GatewayRegister) GetNetworkAdminInfo() string {
	return r.NetworkAdminInfo
}

// GetNetworkAdminInfo gets the network admin ap
func (r *ProviderRegister) GetNetworkAdminInfo() string {
	return r.NetworkAdminInfo
}

// GetRootSigningKey gets the root signing key
func (r *GatewayRegister) GetRootSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.RootSigningKey)
}

// GetRootSigningKey gets the root signing key
func (r *ProviderRegister) GetRootSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.RootSigningKey)
}

// GetSigningKey gets the signing key
func (r *GatewayRegister) GetSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.SigningKey)
}

// GetSigningKey gets the signing key
func (r *ProviderRegister) GetSigningKey() (*fcrcrypto.KeyPair, error) {
	return fcrcrypto.DecodePublicKey(r.SigningKey)
}

// RegisterGateway to register a gateway
func (r *GatewayRegister) RegisterGateway(registerURL string) error {
	url := registerURL + "/registers/gateway"
	err := request.SendJSON(url, r)
	if err != nil {
		return err
	}
	return nil
}

// RegisterProvider to register a provider
func (r *ProviderRegister) RegisterProvider(registerURL string) error {
	url := registerURL + "/registers/provider"
	err := request.SendJSON(url, r)
	if err != nil {
		return err
	}
	return nil
}
