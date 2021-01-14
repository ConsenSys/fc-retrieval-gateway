package messages

// gateway.go contains all messages originting from the gateway

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// GatewaySingleCIDOfferPublishRequest is the request from gateway to provider during start-up asking for cid offers
type GatewaySingleCIDOfferPublishRequest struct {
	MessageType        int32         `json:"message_type"`
	ProtocolVersion    int32         `json:"protocol_version"`
	ProtocolSupported  []int32       `json:"protocol_supported"`
	CIDMin             cid.ContentID `json:"cid_min"`
	CIDMax             cid.ContentID `json:"cid_max"`
	BlockHash          string        `json:"block_hash"`
	TransactionReceipt string        `json:"transaction_receipt"`
	MerkleProof        string        `json:"merkle_proof"`
}

// GatewaySingleCIDOfferPublishResponse is the repsonse to GatewaySingleCIDOfferPublishRequest
type GatewaySingleCIDOfferPublishResponse struct {
	MessageType        int32   `json:"message_type"`
	ProtocolVersion    int32   `json:"protocol_version"`
	ProtocolSupported  []int32 `json:"protocol_supported"`
	PublishedGroupCIDs []struct {
		Nonce      int64         `json:"nonce"`
		ProviderID nodeid.NodeID `json:"provider_id"`
		NumOffers  int64         `json:"num_of_offers"`
		CIDOffers  []struct {
			Price     uint64        `json:"price_per_byte"`
			Expiry    int64         `json:"expiry_date"`
			QoS       uint64        `json:"qos"`
			Signature string        `json:"signature"`
			PieceCID  cid.ContentID `json:"piece_cid"`
		} `json:"cid_offers"`
	} `json:"published_group_cids"` // TODO: Need to check if this is right.
}

// GatewaySingleCIDOfferPublishResponseAck is the acknowledgement to GatewaySingleCIDOfferPublishResponse
type GatewaySingleCIDOfferPublishResponseAck struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	ProtocolSupported []int32 `json:"protocol_supported"`
	CIDOffersAck      []struct {
		Nonce     int64  `json:"nonce"`
		Signature string `json:"signature"`
	} `json:"cid_offers_ack"` // TODO: Need to check if this is right.
}

// GatewayDHTDiscoverRequest is the request from gateway to gateway to discover cid offer
type GatewayDHTDiscoverRequest struct {
	MessageType       int32         `json:"message_type"`
	ProtocolVersion   int32         `json:"protocol_version"`
	ProtocolSupported []int32       `json:"protocol_supported"`
	PieceCID          cid.ContentID `json:"piece_cid"`
	Nonce             int64         `json:"nonce"`
	TTL               int64         `json:"ttl"`
}

// GatewayDHTDiscoverResponse is the response to GatewayDHTDiscoverRequest
type GatewayDHTDiscoverResponse struct {
	MessageType       int32                 `json:"message_type"`
	ProtocolVersion   int32                 `json:"protocol_version"`
	ProtocolSupported []int32               `json:"protocol_supported"`
	PieceCID          cid.ContentID         `json:"piece_cid"`
	Nonce             int64                 `json:"nonce"`
	Found             bool                  `json:"found"`
	Signature         string                `json:"signature"`
	CIDGroupInfo      []CIDGroupInformation `json:"cid_group_information"`
}
