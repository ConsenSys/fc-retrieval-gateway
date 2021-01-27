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

// HandleClientNetworkEstablishment is used to handle initial establishment http request from client
func handleClientNetworkEstablishment(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	g := gateway.GetSingleInstance()

	clientID, challenge, ttl, err := fcrmessages.DecodeClientEstablishmentRequest(request)
	if err != nil {
		s := "Client Establishment: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	logging.Trace("Client Establishment from %s with challenge %s and ttl %d", clientID.ToString(), challenge, ttl)

	now := util.GetTimeImpl().Now().Unix()
	if now > ttl {
		// TODO how to just drop the connection?
		return
	}

	// Construct message
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, challenge)
	if err != nil {
		// TODO for the moment just blow up!
		panic(err)
		// s := "Client Establishment: Internal error signing."
		// logging.Error(s + err.Error())
		// rest.Error(w, s, http.StatusInternalServerError)

	}

	response, err := fcrmessages.EncodeClientEstablishmentResponse(g.GatewayID, challenge, sig)
	if err != nil {
		s := "Client Establishment: Error encoding payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	w.WriteJson(response)
}
