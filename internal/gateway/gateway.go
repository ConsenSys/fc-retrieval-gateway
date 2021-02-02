package gateway

import (
	"errors"
	"net"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/offers"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmerkletrie"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

const (
	protocolVersion   = 1 // Main protocol version
	protocolSupported = 1 // Alternative protocol version
)

// CommunicationChannels holds channels necessary to manage the live tcp connection thread
// CommsLock is used to ensure only one thread can access the tcp connection thread.
// Send true to InterruptRequestChan to indicate the attempt to interrupt tcp connection thread.
// When tcp connection thread is ready, it will send true to interruptResponseChan.
// Send any request in bytes directly to CommsRequestChan.
// The connection thread will send the response to CommsResponseChan if there’s any response expected.
// The connection thread will also send any error to the CommsResponseError or nil if the request is successful.
// Send false to InterruptRequestChan to indicate the attempt to end the interruption.

// CommunicationChannel holds the connection for sending outgoing TCP requests.
// CommsLock is used to ensure only one thread can access the tcp connection at any time.
// Conn is the net connection for sending outgoing TCP requests.
type CommunicationChannel struct {
	CommsLock sync.RWMutex
	Conn      net.Conn
}

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

	// GatewayAddressMap stores mapping from gateway id (big int in string repr) to its address.
	GatewayAddressMap     map[string](string)
	GatewayAddressMapLock sync.RWMutex

	// ProviderAddressMap stores mapping from provider id (big int in string repr) to its address.
	ProviderAddressMap     map[string](string)
	ProviderAddressMapLock sync.RWMutex

	// ActiveGateways store connected active gateways for outgoing request:
	// A map from gateway id (big int in string repr) to a CommunicationChannel.
	ActiveGateways     map[string](*CommunicationChannel)
	ActiveGatewaysLock sync.RWMutex

	// ActiveProviders store connected active providers for outgoing request:
	// A map from provider id (big in in string repr) to a CommunicationChannel.
	ActiveProviders     map[string](*CommunicationChannel)
	ActiveProvidersLock sync.RWMutex

	// Offers, it is threadsafe.
	Offers *offers.Offers

	// RegistrationBlockHash is the hash of the block that registers this gateway
	// RegistrationTransactionReceipt is the transaction receipt containing the registration event
	// RegistrationMerkleRoot is the root of the merkle trie containing the transaction receipt
	// RegistrationMerkleProof proves the transaction receipt is part of the block
	RegistrationBlockHash          string
	RegistrationTransactionReceipt string
	RegistrationMerkleRoot         string
	RegistrationMerkleProof        *fcrmerkletrie.FCRMerkleProof
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

		gatewayPrivateKey, err := fcrcrypto.DecodePrivateKey(conf.GatewayPrivKey)
		if err != nil {
			logging.ErrorAndPanic("Error decoding Gateway Private Key: %s", err)
		}
		gatewayID, err := nodeid.NewNodeIDFromString(conf.GatewayID)
		if err != nil {
			logging.ErrorAndPanic("Error decoding node id: %s", err)
		}

		gatewayPrivateKeyVersion := fcrcrypto.DecodeKeyVersion(conf.GatewayKeyVersion)

		instance = &Gateway{
			ProtocolVersion:          protocolVersion,
			ProtocolSupported:        []int32{protocolVersion, protocolSupported},
			GatewayAddressMap:        make(map[string](string)),
			GatewayAddressMapLock:    sync.RWMutex{},
			ProviderAddressMap:       make(map[string](string)),
			ProviderAddressMapLock:   sync.RWMutex{},
			ActiveGateways:           make(map[string](*CommunicationChannel)),
			ActiveGatewaysLock:       sync.RWMutex{},
			ActiveProviders:          make(map[string](*CommunicationChannel)),
			ActiveProvidersLock:      sync.RWMutex{},
			GatewayPrivateKey:        gatewayPrivateKey,
			GatewayPrivateKeyVersion: gatewayPrivateKeyVersion,
			GatewayID:                gatewayID,
			Offers:                   offers.GetSingleInstance(),
			RegistrationBlockHash:    "TODO",
		}
	})
	return instance
}

// RegisterGatewayCommunication registers a gateway communication
func RegisterGatewayCommunication(id *nodeid.NodeID, gComm *CommunicationChannel) error {
	if instance == nil {
		return errors.New("Error: instance not created")
	}
	instance.ActiveGatewaysLock.Lock()
	defer instance.ActiveGatewaysLock.Unlock()
	_, exist := instance.ActiveGateways[id.ToString()]
	if exist {
		return errors.New("Error: connection existed")
	}
	instance.ActiveGateways[id.ToString()] = gComm
	return nil
}

// DeregisterGatewayCommunication deregisters a gateway communication
// Fail silently
func DeregisterGatewayCommunication(id *nodeid.NodeID) {
	if instance != nil {
		instance.ActiveGatewaysLock.Lock()
		defer instance.ActiveGatewaysLock.Unlock()
		_, exist := instance.ActiveGateways[id.ToString()]
		if exist {
			delete(instance.ActiveGateways, id.ToString())
		}
	}
}

// RegisterProviderCommunication registers a provider communication
func RegisterProviderCommunication(id *nodeid.NodeID, pComm *CommunicationChannel) error {
	if instance == nil {
		return errors.New("Error: instance not created")
	}
	instance.ActiveProvidersLock.Lock()
	defer instance.ActiveProvidersLock.Unlock()
	_, exist := instance.ActiveProviders[id.ToString()]
	if exist {
		return errors.New("Error: connection existed")
	}
	instance.ActiveProviders[id.ToString()] = pComm
	return nil
}

// DeregisterProviderCommunication deregisters a provider communication
func DeregisterProviderCommunication(id *nodeid.NodeID) error {
	if instance == nil {
		return errors.New("Error: instance not created")
	}
	instance.ActiveProvidersLock.Lock()
	defer instance.ActiveProvidersLock.Unlock()
	_, exist := instance.ActiveProviders[id.ToString()]
	if !exist {
		return errors.New("Error: connection not existed")
	}
	delete(instance.ActiveProviders, id.ToString())
	return nil
}
