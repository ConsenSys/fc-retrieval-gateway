package messages

import "math/big"


// ClientEstablishmentRequest is the request from client to gateway to establish connection
type ClientEstablishmentRequest struct {
	MessageType       int32     `json:"message_type"`
	ProtocolVersion   int32     `json:"protocol_version"`
	ProtocolSupported []int32   `json:"protocol_supported"`
	Challenge         string  `json:"challenge"`
	TTL               int64     `json:"ttl"`
}

// ClientEstablishmentResponse is the response to ClientEstablishmentRequest
type ClientEstablishmentResponse struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	Challenge    string              `json:"challenge"`
	Signature    string              `json:"signature"`
}

// ClientStandardDiscoverRequest is the requset from client to gateway to ask for cid offer
type ClientStandardDiscoverRequest struct {
	CommonFields CommonRequestMessageFields `json:"common_fields"`
	PieceCID     big.Int             `json:"piece_cid"`
	Nonce        int                 `json:"nonce"`
	TTL          string              `json:"ttl"`
}

// ClientStandardDiscoverResponse is the response to ClientStandardDiscoverResponse
type ClientStandardDiscoverResponse struct {
	CommonFields CommonRequestMessageFields   `json:"common_fields"`
	PieceCID     big.Int               `json:"piece_cid"`
	Nonce        int                   `json:"nonce"`
	Found        bool                  `json:"found"`
	Signature    string                `json:"signature"`
	CIDGroupInfo []CIDGroupInformation `json:"cid_group_information"`
}

// ClientDHTDiscoverRequest is the request from client to gateway to ask for cid offer using DHT
type ClientDHTDiscoverRequest struct {
	CommonFields       CommonRequestMessageFields `json:"common_fields"`
	PieceCID           big.Int             `json:"piece_cid"`
	Nonce              int                 `json:"nonce"`
	TTL                string              `json:"ttl"`
	NumDHT             int                 `json:"num_dht"`
	IncrementalResults bool                `json:"incremental_results"`
}

// ClientDHTDiscoverResponse is the response to ClientDHTDiscoverRequest
type ClientDHTDiscoverResponse struct {
	CommonFields  CommonRequestMessageFields              `json:"common_fields"`
	Contacted     []ClientStandardDiscoverResponse `json:"contacted_gateways"`
	UnContactable []struct {
		GatewayID big.Int `json:"gateway_id"`
		Nonce     int     `json:"nonce"`
	} `json:"uncontactable_gateways"`
}

