package clientapi

// Copyright (C) 2020 ConsenSys Software Inc
import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/api/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientDHTCIDDiscover is used to handle client request for cid offer
func (c *ClientAPI) HandleClientDHTCIDDiscover(w rest.ResponseWriter, content []byte) {
	request := messages.ClientDHTDiscoverRequest{}
	err := json.Unmarshal(content, &request)
	if err != nil {
		s := "Client DHT CID Discovery: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	// First check if the message can be discarded
	if time.Now().Unix() > request.TTL {
		// Message expired.
		return
	}
	// Use DHT to get response.
	c.gateway.GatewayAddressMapLock.RLock()
	defer c.gateway.GatewayAddressMapLock.RUnlock()

	if len(c.gateway.GatewayAddressMap) < int(request.NumDHT) {
		s := "Gateway does not store enough peers."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}
	gatewayIDs := make([]*nodeid.NodeID, int(request.NumDHT))

	// TODO: Need to add an algorithm to select gateways from the map.
	// For now, it is random.
	i := 0
	for k := range c.gateway.GatewayAddressMap {
		if i >= int(request.NumDHT) {
			break
		}
		gatewayIDs[i], _ = nodeid.NewNodeIDFromString(k)
		i++
	}
	// Construct response
	// TODO: Right now, it ignores the incremental result filed.
	// Will return all in one message.
	response := messages.ClientDHTDiscoverResponse{
		MessageType:       messages.ClientDHTDiscoverResponseType,
		ProtocolVersion:   c.gateway.ProtocolVersion,
		ProtocolSupported: c.gateway.ProtocolSupported,
		Contacted:         make([]messages.GatewayDHTDiscoverResponse, 0),
		UnContactable:     make([]nodeid.NodeID, 0)}
	// Now requesting gateways.
	for _, id := range gatewayIDs {
		res, err := gatewayapi.RequestGatewayDHTDiscover(&request.PieceCID, id, c.gateway)
		if err != nil {
			response.UnContactable = append(response.UnContactable, *id)
		} else {
			response.Contacted = append(response.Contacted, *res)
		}
	}
	response.Nonce = request.Nonce
	w.WriteJson(response)
}
