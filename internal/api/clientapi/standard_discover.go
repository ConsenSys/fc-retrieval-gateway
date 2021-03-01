package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmerkletree"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
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

	roots := make([]string, 0)
	proofs := make([]fcrmerkletree.FCRMerkleProof, 0)
	fundedPaymentChannel := make([]bool, 0)

	for _, offer := range offers {
		tree := offer.GetMerkleTrie()
		roots = append(roots, tree.GetMerkleRoot())
		proof, err := tree.GenerateMerkleProof(pieceCID)
		if err != nil {
			s := "Internal error: Error generating proof."
			logging.Error(s + err.Error())
			rest.Error(w, s, http.StatusBadRequest)
			return
		}
		proofs = append(proofs, *proof)
		fundedPaymentChannel = append(fundedPaymentChannel, false) // TODO, Need to find a way to check if having payment channel set up for a given provider.
	}

	// Construct response
	response, err := fcrmessages.EncodeClientStandardDiscoverResponse(pieceCID, nonce, exists, offers, roots, proofs, fundedPaymentChannel)
	if err != nil {
		s := "Internal error: Error encoding payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// Sign the message
	if response.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, msg)
	}) != nil {
		s := "Internal error."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	w.WriteJson(response)
}
