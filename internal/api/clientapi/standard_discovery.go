package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientStandardCIDDiscover is used to handle client request for cid offer
func (c *ClientAPI) HandleClientStandardCIDDiscover(w rest.ResponseWriter, content []byte) {
	request := messages.ClientStandardDiscoverRequest{}
	err := json.Unmarshal(content, &request)
	if err != nil {
		s := "Client Standard CID Discovery: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	// First check if the message can be discarded
	if time.Now().Unix() > request.TTL {
		// Message expired.
		return
	}
	// Respond
	offers, exists := c.gateway.Offers.GetOffers(&request.PieceCID)

	// Construct response
	response := messages.ClientStandardDiscoverResponse{
		MessageType:       messages.GatewayDHTDiscoverResponseType,
		ProtocolVersion:   c.gateway.ProtocolVersion,
		ProtocolSupported: c.gateway.ProtocolSupported,
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
	sig, err := fcrcrypto.SignMessage(c.gateway.GatewayPrivateKey, c.gateway.GatewayPrivateKeyVersion, response)
	if err != nil {
		s := "Internal error."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	response.Signature = sig
	w.WriteJson(response)
}
