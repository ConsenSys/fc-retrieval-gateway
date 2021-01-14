package gatewayapi

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleGatewayDHTDiscoverRequest(conn net.Conn, request *messages.GatewayDHTDiscoverRequest) error {
	// First check if the message can be discarded.
	if time.Now().Unix() > request.TTL {
		// Message discarded.
		return nil
	}
	// Respond to DHT CID Discover Request
	// Get gateway core struct
	g := gateway.GetSingleInstance()
	offers, exists := g.Offers.GetOffers(&request.PieceCID)

	// Construct response
	response := messages.GatewayDHTDiscoverResponse{
		MessageType:       messages.GatewayDHTDiscoverResponseType,
		ProtocolVersion:   g.ProtocolVersion,
		ProtocolSupported: g.ProtocolSupported,
		PieceCID:          request.PieceCID,
		Nonce:             request.Nonce}

	if exists {
		response.Found = true
		response.CIDGroupInfo = make([]messages.CIDGroupInformation, len(offers))
		for i, offer := range offers {
			response.CIDGroupInfo[i].ProviderID = *offer.NodeID
			response.CIDGroupInfo[i].Price = offer.Price
			response.CIDGroupInfo[i].Expiry = offer.Expiry
			response.CIDGroupInfo[i].QoS = offer.QoS
			response.CIDGroupInfo[i].Signature = offer.Signature
			response.CIDGroupInfo[i].MerkleProof = offer.MerkleProof
			response.CIDGroupInfo[i].FundedPaymentChannel = offer.FundedPaymentChannel
		}
	} else {
		response.Found = false
		response.CIDGroupInfo = make([]messages.CIDGroupInformation, 0)
	}
	// Sign the message
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, response)
	if err != nil {
		// Ignored.
		logging.Error("Error in signing message.")
	}
	response.Signature = sig
	// Send message
	data, _ := json.Marshal(response)
	return tcpcomms.SendTCPMessage(conn, messages.GatewayDHTDiscoverResponseType, data, settings.DefaultTCPInactivityTimeout)
}

// RequestGatewayDHTDiscover is used to request a DHT CID Discover
func RequestGatewayDHTDiscover(cid *cid.ContentID, gatewayID *nodeid.NodeID, g *gateway.Gateway) (*messages.GatewayDHTDiscoverResponse, error) {
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
	request := messages.GatewayDHTDiscoverRequest{
		MessageType:       messages.GatewayDHTDiscoverRequestType,
		ProtocolVersion:   g.ProtocolVersion,
		ProtocolSupported: g.ProtocolSupported,
		PieceCID:          *cid,
		Nonce:             1,                                       // TODO, Add nonce
		TTL:               time.Now().Add(10 * time.Second).Unix(), // TODO, ADD TTL, for now 10 seconds
	}
	err = tcpcomms.SendMessageWithType(pComm.Conn, messages.GatewayDHTDiscoverRequestType, &request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	// Get a response
	msgType, data, err := tcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultLongTCPInactivityTimeout)
	if err != nil && tcpcomms.IsTimeoutError(err) {
		// Timeout can be ignored. Since this message can expire.
		return nil, nil
	} else if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterGatewayCommunication(gatewayID)
		return nil, err
	}
	if msgType == messages.GatewayDHTDiscoverResponseType {
		response := messages.GatewayDHTDiscoverResponse{}
		if json.Unmarshal(data, &response) == nil {
			// Message is valid.
			return &response, nil
		}
	}
	// Message is invalid.
	return nil, errors.New("invalid message")
}
