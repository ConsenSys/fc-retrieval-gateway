package gatewayapi

import (
	"errors"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleGatewayDHTDiscoverRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	gatewayID, pieceCID, nonce, ttl, _, _, err := fcrmessages.DecodeGatewayDHTDiscoverRequest(request)
	if err != nil {
		// Reply with invalid message
		return fcrtcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
	}
	// Get the gateway's signing key
	g.RegisteredGatewaysMapLock.RLock()
	defer g.RegisteredGatewaysMapLock.RUnlock()
	_, ok := g.RegisteredGatewaysMap[gatewayID.ToString()]
	if !ok {
		return errors.New("Gateway public key not found")
	}
	pubKey, err := g.RegisteredProvidersMap[gatewayID.ToString()].GetSigningKey()
	if err != nil {
		return err
	}
	// First verify the message
	if request.Verify(pubKey) != nil {
		return errors.New("Fail to verify the request")
	}

	// Second check if the message can be discarded.
	if time.Now().Unix() > ttl {
		// Message discarded.
		return nil
	}
	// Respond to DHT CID Discover Request
	offers, exists := g.Offers.GetOffers(pieceCID)

	suboffers := make([]cidoffer.SubCIDOffer, 0)
	fundedPaymentChannel := make([]bool, 0)

	for _, offer := range offers {
		suboffer, err := offer.GenerateSubCIDOffer(pieceCID)
		if err != nil {
			return err
		}
		suboffers = append(suboffers, *suboffer)
		fundedPaymentChannel = append(fundedPaymentChannel, false) // TODO, Need to find a way to check if having payment channel set up for a given provider.
	}

	// Construct message
	response, err := fcrmessages.EncodeGatewayDHTDiscoverResponse(pieceCID, nonce, exists, suboffers, fundedPaymentChannel)
	if err != nil {
		return err
	}

	// Sign the response
	if response.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		return errors.New("Error in signing the response")
	}
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}

// RequestGatewayDHTDiscover is used to request a DHT CID Discover
func RequestGatewayDHTDiscover(cid *cid.ContentID, gatewayID *nodeid.NodeID) (*fcrmessages.FCRMessage, error) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// Get the connection to the gateway.
	pComm, err := g.GatewayCommPool.GetConnForRequestingNode(gatewayID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	request, err := fcrmessages.EncodeGatewayDHTDiscoverRequest(g.GatewayID, cid, 1, time.Now().Add(10*time.Second).Unix(), "", "") // TODO, ADD nonce and TTL
	if err != nil {
		return nil, err
	}
	// Sign the request
	if request.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		return nil, errors.New("Error in signing the request")
	}
	// Send the request
	err = fcrtcpcomms.SendTCPMessage(pComm.Conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		g.GatewayCommPool.DeregisterNodeCommunication(gatewayID)
		return nil, err
	}
	// Get a response
	response, err := fcrtcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultTCPInactivityTimeout)
	if err != nil && fcrtcpcomms.IsTimeoutError(err) {
		// Timeout can be ignored. Since this message can expire.
		return nil, nil
	} else if err != nil {
		g.GatewayCommPool.DeregisterNodeCommunication(gatewayID)
		return nil, err
	}
	// Verify the response
	// Get the gateway's signing key
	g.RegisteredGatewaysMapLock.RLock()
	defer g.RegisteredGatewaysMapLock.RUnlock()
	pubKey, err := g.RegisteredProvidersMap[gatewayID.ToString()].GetSigningKey()
	if err != nil {
		return nil, err
	}
	if response.Verify(pubKey) != nil {
		return nil, errors.New("Fail to verify the response")
	}
	return response, nil
}
