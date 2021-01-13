package gatewayapi

import (
	"encoding/json"
	"errors"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/providerapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// RequestSingleCIDOffers is used at start-up to request a set of single CID Offers
// from a provider with a given provider id.
func RequestSingleCIDOffers(cidMin, cidMax cid.ContentID, providerID *nodeid.NodeID, g *gateway.Gateway) (*messages.GatewaySingleCIDOfferPublishResponse, error) {
	// Get the connection to provider.
	pComm, err := providerapi.GetConnForRequestingProvider(providerID, g)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterProviderCommunication(providerID)
		return nil, err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	request := messages.GatewaySingleCIDOfferPublishRequest{
		MessageType:        messages.GatewaySingleCIDOfferPublishRequestType,
		ProtocolVersion:    g.ProtocolVersion,
		ProtocolSupported:  g.ProtocolSupported,
		CIDMin:             cidMin,
		CIDMax:             cidMax,
		BlockHash:          g.RegistrationBlockHash,
		TransactionReceipt: g.RegistrationTransactionReceipt,
		MerkleProof:        g.RegistrationMerkleProof}
	err = tcpcomms.SendMessageWithType(pComm.Conn, messages.GatewayDHTDiscoverRequestType, &request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterProviderCommunication(providerID)
		return nil, err
	}
	// Get a response.
	msgType, data, err := tcpcomms.ReadTCPMessage(pComm.Conn, settings.DefaultLongTCPInactivityTimeout)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterProviderCommunication(providerID)
		return nil, err
	}
	if msgType == messages.GatewaySingleCIDOfferPublishResponseType {
		response := messages.GatewaySingleCIDOfferPublishResponse{}
		if json.Unmarshal(data, &response) == nil {
			// Message is valid.
			return &response, nil
		}
	}
	// Message is invalid.
	return nil, errors.New("invalid message")
}

// AcknowledgeSingleCIDOffers is used to acknowledge a response
func AcknowledgeSingleCIDOffers(response *messages.GatewaySingleCIDOfferPublishResponse, providerID *nodeid.NodeID, g *gateway.Gateway) error {
	// Get the connection to provider.
	pComm, err := providerapi.GetConnForRequestingProvider(providerID, g)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterProviderCommunication(providerID)
		return err
	}
	pComm.CommsLock.Lock()
	defer pComm.CommsLock.Unlock()
	// Construct message
	cidOffersAck := make([]struct {
		Nonce     int64  `json:"nonce"`
		Signature string `json:"signature"`
	}, len(response.PublishedGroupCIDs))
	for i := 0; i < len(response.PublishedGroupCIDs); i++ {
		cidOffersAck[i].Nonce = response.PublishedGroupCIDs[i].Nonce
		// Sign the offer
		sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, response.PublishedGroupCIDs[i])
		if err != nil {
			// Ignored.
			logging.Error("Error in signing message.")
		}
		cidOffersAck[i].Signature = sig
	}
	request := messages.GatewaySingleCIDOfferPublishResponseAck{
		MessageType:       messages.GatewaySingleCIDOfferPublishResponseAckType,
		ProtocolVersion:   g.ProtocolVersion,
		ProtocolSupported: g.ProtocolSupported,
		CIDOffersAck:      cidOffersAck}
	err = tcpcomms.SendMessageWithType(pComm.Conn, messages.GatewaySingleCIDOfferPublishResponseAckType, &request, settings.DefaultTCPInactivityTimeout)
	if err != nil {
		pComm.Conn.Close()
		gateway.DeregisterProviderCommunication(providerID)
		return err
	}
	// ACK Success
	return nil
}
