package gatewayapi

import (
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

func handleGatewayDHTDiscoverRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	pieceCID, nonce, ttl, err := fcrmessages.DecodeGatewayDHTDiscoverRequest(request)
	if err != nil {
		// Reply with invalid message
		return fcrtcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
	}

	// First check if the message can be discarded.
	if time.Now().Unix() > ttl {
		// Message discarded.
		return nil
	}
	// Respond to DHT CID Discover Request
	// Get gateway core struct
	g := gateway.GetSingleInstance()
	offers, exists := g.Offers.GetOffers(pieceCID)

	// Construct message
	response, err := fcrmessages.EncodeGatewayDHTDiscoverResponse(pieceCID, nonce, exists, offers)
	if err != nil {
		return err
	}

	// Sign the message
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, response)
	if err != nil {
		// Ignored.
		logging.Error("Error in signing message.")
	}
	response.SetSignature(sig)
	// Send message
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}

// RequestGatewayDHTDiscover is used to request a DHT CID Discover
func RequestGatewayDHTDiscover(cid *cid.ContentID, gatewayID *nodeid.NodeID) (*fcrmessages.FCRMessage, error) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// Get the connection to the gateway.
	pComm, err := GetConnForRequestingGateway(gatewayID, g)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	request, err := fcrmessages.EncodeGatewayDHTDiscoverRequest(cid, 1, time.Now().Add(10*time.Second).Unix()) // TODO, ADD nonce and TTL
	if err != nil {
		return nil, err
	}
	err = fcrtcpcomms.SendTCPMessage(pComm.Conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	// Get a response
	response, err := fcrtcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultTCPInactivityTimeout)
	if err != nil && fcrtcpcomms.IsTimeoutError(err) {
		// Timeout can be ignored. Since this message can expire.
		return nil, nil
	} else if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	return response, nil
}
