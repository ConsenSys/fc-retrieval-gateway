package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"net/http"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientStandardCIDDiscover is used to handle client request for cid offer
func handleClientStandardCIDDiscover(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	g := gateway.GetSingleInstance()

	pieceCID, nonce, ttl, err := fcrmessages.DecodeClientStandardDiscoverRequest(request)
	if err != nil {
		s := "Client Standard CID Discovery: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	now := util.GetTimeImpl().Now().Unix()
	if now > ttl {
		// Drop the connection
		return
	}

	// Search for offesr.
	offers, exists := g.Offers.GetOffers(pieceCID)

	// Construct response
	response, err := fcrmessages.EncodeClientStandardDiscoverResponse(pieceCID, nonce, exists, offers)
	if err != nil {
		s := "Internal error: Error encoding payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// Sign the message
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, response)
	if err != nil {
		s := "Internal error."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	// Set signature
	response.SetSignature(sig)
	w.WriteJson(response)
}
