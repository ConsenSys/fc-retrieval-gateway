package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"net/http"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientDHTCIDDiscover is used to handle client request for cid offer
func handleClientDHTCIDDiscover(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	g := gateway.GetSingleInstance()

	cid, nonce, ttl, numDHT, _, _, _, err := fcrmessages.DecodeClientDHTDiscoverRequest(request)
	if err != nil {
		s := "Client DHT CID Discovery: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// First check if the message can be discarded
	if time.Now().Unix() > ttl {
		// Message expired.
		return
	}
	// Use DHT to get response.
	g.RegisteredGatewaysMapLock.RLock()
	defer g.RegisteredGatewaysMapLock.RUnlock()

	if len(g.RegisteredGatewaysMap) < int(numDHT) {
		s := "Gateway does not store enough peers."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	gatewayIDs := make([]*nodeid.NodeID, int(numDHT))

	// TODO: Need to add an algorithm to select gateways from the map.
	// For now, it is random.
	i := 0
	for k := range g.RegisteredGatewaysMap {
		if i >= int(numDHT) {
			break
		}
		gatewayIDs[i], _ = nodeid.NewNodeIDFromHexString(k)
		i++
	}
	// Construct response
	// TODO: Right now, it ignores the incremental result filed.
	// Will return all in one message.
	// Now requesting gateways.
	contacted := make([]fcrmessages.FCRMessage, 0)
	unContactable := make([]nodeid.NodeID, 0)
	for _, id := range gatewayIDs {
		res, err := gatewayapi.RequestGatewayDHTDiscover(cid, id)
		if err != nil {
			unContactable = append(unContactable, *id)
		} else {
			contacted = append(contacted, *res)
		}
	}

	response, err := fcrmessages.EncodeClientDHTDiscoverResponse(contacted, unContactable, nonce)
	if err != nil {
		s := "Internal error: Fail to encode response."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}

	// Sign the message
	// Sign message
	if response.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		s := "Internal error."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	w.WriteJson(response)
}
