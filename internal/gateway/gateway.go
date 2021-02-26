package gateway

import (
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/offers"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmerkletree"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

const (
	protocolVersion   = 1 // Main protocol version
	protocolSupported = 1 // Alternative protocol version
)

// Gateway holds the main data structure for the whole gateway.
type Gateway struct {
	ProtocolVersion   int32
	ProtocolSupported []int32

	// GatewayID of this gateway
	GatewayID *nodeid.NodeID

	// Gateway Private Key and algorithm of this gateway
	GatewayPrivateKey *fcrcrypto.KeyPair

	// GatewayPrivateKeyVersion is the key version number of the private key.
	GatewayPrivateKeyVersion *fcrcrypto.KeyVersion

	// RegisteredGatewaysMap stores mapping from gateway id (big int in string repr) to its registration info
	RegisteredGatewaysMap     map[string]register.RegisteredNode
	RegisteredGatewaysMapLock sync.RWMutex

	// RegisteredProvidersMap stores mapping from provider id (big int in string repr) to its registration info
	RegisteredProvidersMap     map[string]register.RegisteredNode
	RegisteredProvidersMapLock sync.RWMutex

	// GatewayCommPool manages connection for outgoing request to gateways
	GatewayCommPool *fcrtcpcomms.CommunicationPool

	// ProviderCommPool manages connection for outgoing request to providers
	ProviderCommPool *fcrtcpcomms.CommunicationPool

	// Offers, it is threadsafe.
	Offers *offers.Offers

	// RegistrationBlockHash is the hash of the block that registers this gateway
	// RegistrationTransactionReceipt is the transaction receipt containing the registration event
	// RegistrationMerkleRoot is the root of the merkle trie containing the transaction receipt
	// RegistrationMerkleProof proves the transaction receipt is part of the block
	RegistrationBlockHash          string
	RegistrationTransactionReceipt string
	RegistrationMerkleRoot         string
	RegistrationMerkleProof        *fcrmerkletree.FCRMerkleProof
}

// Single instance of the gateway
var instance *Gateway
var doOnce sync.Once

// GetSingleInstance returns the single instance of the gateway
func GetSingleInstance(confs ...*settings.AppSettings) *Gateway {
	doOnce.Do(func() {
		if len(confs) == 0 {
			logging.ErrorAndPanic("No settings supplied to Gateway start-up")
		}
		if len(confs) != 1 {
			logging.ErrorAndPanic("More than one sets of settings supplied to Gateway start-up")
		}
		conf := confs[0]

		gatewayID, err := nodeid.NewNodeIDFromString(conf.GatewayID)
		if err != nil {
			logging.ErrorAndPanic("Error decoding node id: %s", err)
		}

		instance = &Gateway{
			ProtocolVersion:          protocolVersion,
			ProtocolSupported:        []int32{protocolVersion, protocolSupported},
			RegisteredGatewaysMap:          make(map[string]register.RegisteredNode),
			RegisteredGatewaysMapLock:      sync.RWMutex{},
			RegisteredProvidersMap:         make(map[string]register.RegisteredNode),
			RegisteredProvidersMapLock:     sync.RWMutex{},
			GatewayPrivateKey:        nil,
			GatewayPrivateKeyVersion: nil,
			GatewayID:                gatewayID,
			Offers:                   offers.GetSingleInstance(),
			RegistrationBlockHash:          "TODO",
			RegistrationTransactionReceipt: "TODO",
			RegistrationMerkleRoot:         "TODO",
			RegistrationMerkleProof:        nil, //TODO
		}
		instance.GatewayCommPool = fcrtcpcomms.NewCommunicationPool(instance.RegisteredGatewaysMap, &instance.RegisteredGatewaysMapLock)
		instance.ProviderCommPool = fcrtcpcomms.NewCommunicationPool(instance.RegisteredProvidersMap, &instance.RegisteredProvidersMapLock)
	})
	return instance
}
