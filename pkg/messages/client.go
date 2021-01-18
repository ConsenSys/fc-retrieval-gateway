package messages

// client.go contains all the messages originating from the client

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ClientCommonRequestFields common fields in requests from client to gateway.
type ClientCommonRequestFields struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	ProtocolSupported []int32 `json:"protocol_supported"`
	ClientID          string  `json:"client_id"`
	TTL               int64   `json:"ttl"`
	Signature         string  `json:"signature"`
}

// Set sets the common fields.
func (c *ClientCommonRequestFields) Set(
	messageType int32, protocolVersion int32, protocolSupported []int32,
	clientID string, ttl int64) {
		c.MessageType = messageType
		c.ProtocolVersion = protocolVersion
		c.ProtocolSupported = protocolSupported
		c.ClientID = clientID
		c.TTL = ttl
	}

// SetSignature sets the signature field.
func (c *ClientCommonRequestFields) SetSignature(signature string) {
	c.Signature = signature
}


// ClientEstablishmentRequest is the request from client to gateway to establish connection
type ClientEstablishmentRequest struct {
	Challenge         string  `json:"challenge"`
	ClientCommonRequestFields
}

// ClientEstablishmentResponse is the response to ClientEstablishmentRequest
type ClientEstablishmentResponse struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	ProtocolSupported []int32 `json:"protocol_supported"`
	GatewayID         string  `json:"gateway_id"`
	Challenge         string  `json:"challenge"`
	Signature         string  `json:"signature"`
}

// ClientStandardDiscoverRequest is the requset from client to gateway to ask for cid offer
type ClientStandardDiscoverRequest struct {
	PieceCID          cid.ContentID `json:"piece_cid"`
	Nonce             int64         `json:"nonce"`
	ClientCommonRequestFields
}

// ClientStandardDiscoverResponse is the response to ClientStandardDiscoverResponse
type ClientStandardDiscoverResponse struct {
	MessageType       int32                 `json:"message_type"`
	ProtocolVersion   int32                 `json:"protocol_version"`
	ProtocolSupported []int32               `json:"protocol_supported"`
	PieceCID          cid.ContentID         `json:"piece_cid"`
	Nonce             int64                 `json:"nonce"`
	Found             bool                  `json:"found"`
	Signature         string                `json:"signature"`
	CIDGroupInfo      []CIDGroupInformation `json:"cid_group_information"`
}

// ClientDHTDiscoverRequest is the request from client to gateway to ask for cid offer using DHT
type ClientDHTDiscoverRequest struct {
	PieceCID           cid.ContentID `json:"piece_cid"`
	Nonce              int64         `json:"nonce"`
	NumDHT             int64         `json:"num_dht"`
	IncrementalResults bool          `json:"incremental_results"`
	ClientCommonRequestFields
}

// ClientDHTDiscoverResponse is the response to ClientDHTDiscoverRequest
type ClientDHTDiscoverResponse struct {
	MessageType       int32                        `json:"message_type"`
	ProtocolVersion   int32                        `json:"protocol_version"`
	ProtocolSupported []int32                      `json:"protocol_supported"`
	Contacted         []GatewayDHTDiscoverResponse `json:"contacted_gateways"`
	UnContactable     []nodeid.NodeID              `json:"uncontactable_gateways"`
	Nonce             int64                        `json:"nonce"`
}

// ClientCIDGroupPublishDHTAckRequest is the request from client to provider to request the signed ack of a cid group publish
type ClientCIDGroupPublishDHTAckRequest struct {
	MessageType       int32         `json:"message_type"`
	ProtocolVersion   int32         `json:"protocol_version"`
	ProtocolSupported []int32       `json:"protocol_supported"`
	PieceCID          cid.ContentID `json:"piece_cid"`
	GatewayID         nodeid.NodeID `json:"gateway_id"`
}

// ClientCIDGroupPublishDHTAckResponse is the response to ClientCIDGroupPublishDHTAckRequest
type ClientCIDGroupPublishDHTAckResponse struct {
	MessageType             int32                             `json:"message_type"`
	ProtocolVersion         int32                             `json:"protocol_version"`
	ProtocolSupported       []int32                           `json:"protocol_supported"`
	PieceCID                cid.ContentID                     `json:"piece_cid"`
	GatewayID               nodeid.NodeID                     `json:"gateway_id"`
	Found                   bool                              `json:"found"`
	CIDGroupPublishToDHT    ProviderDHTPublishGroupCIDRequest `json:"cid_group_publish_to_dht"`
	CIDGroupPublishToDHTAck ProviderDHTPublishGroupCIDAck     `json:"cid_group_publish_to_dht_ack"`
}
