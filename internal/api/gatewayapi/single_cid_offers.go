package gatewayapi

import (
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// RequestSingleCIDOffers is used at start-up to request a set of single CID Offers
// from a provider with a given provider id.
func RequestSingleCIDOffers(cidMin, cidMax *cid.ContentID, providerID *nodeid.NodeID) (*fcrmessages.FCRMessage, error) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// Get the connection to provider.
	pComm, err := g.ProviderCommPool.GetConnForRequestingNode(providerID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	request, err := fcrmessages.EncodeGatewaySingleCIDOfferPublishRequest(
		g.GatewayID,
		cidMin,
		cidMax,
		g.RegistrationBlockHash,
		g.RegistrationTransactionReceipt,
		g.RegistrationMerkleRoot,
		g.RegistrationMerkleProof,
	)
	if err != nil {
		return nil, err
	}
	err = fcrtcpcomms.SendTCPMessage(pComm.Conn, request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	// Get a response.
	response, err := fcrtcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultLongTCPInactivityTimeout)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	return response, nil
}

// AcknowledgeSingleCIDOffers is used to acknowledge a response
func AcknowledgeSingleCIDOffers(response *fcrmessages.FCRMessage, providerID *nodeid.NodeID) ([]fcrmessages.FCRMessage, error) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// Get the connection to provider.
	pComm, err := g.ProviderCommPool.GetConnForRequestingNode(providerID, fcrtcpcomms.AccessFromGateway)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Decode the response
	cidOffers, err := fcrmessages.DecodeGatewaySingleCIDOfferPublishResponse(response)
	if err != nil {
		return nil, err
	}
	// Construct the message
	cidOfferAcks := make([]fcrmessages.FCRMessage, 0)
	for _, cidOffer := range cidOffers {
		nonce, _, _, err := fcrmessages.DecodeProviderDHTPublishGroupCIDRequest(&cidOffer)
		if err != nil {
			return nil, err
		}
		// Sign the offer
		sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, cidOffer)
		if err != nil {
			return nil, err
		}
		cidOfferAck, err := fcrmessages.EncodeProviderDHTPublishGroupCIDAck(nonce, sig)
		if err != nil {
			return nil, err
		}
		cidOfferAcks = append(cidOfferAcks, *cidOfferAck)
	}
	ack, err := fcrmessages.EncodeGatewaySingleCIDOfferPublishResponseAck(cidOfferAcks)
	if err != nil {
		return nil, err
	}
	// Send ack
	err = fcrtcpcomms.SendTCPMessage(pComm.Conn, ack, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		g.ProviderCommPool.DeregisterNodeCommunication(providerID)
		return nil, err
	}
	// Ack success
	return cidOffers, nil
}
