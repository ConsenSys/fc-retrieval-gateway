package messages

import "math/big"

// CommonRequestMessageFields are shared fileds of all http requests
type CommonRequestMessageFields struct {
	MessageType       int   `json:"message_type"`
	ProtocolVersion   int   `json:"protocol_version"`
	ProtocolSupported []int `json:"protocol_supported"`
}

// ProtocolChangeResponse message is sent to indicate that the gateway is requesting the
// other entity change protocol version
type ProtocolChangeResponse struct {
	MessageType       int   `json:"message_type"`
	DesiredVersion    int   `json:"desired_version"`
}


// ProtocolMismatchResponse message is sent to indicate that there are no common protocol 
// versions between the gateway and the requesting entity.
type ProtocolMismatchResponse struct {
	MessageType       int   `json:"message_type"`
}


// CIDGroupInformation represents a cid group information
type CIDGroupInformation struct {
	ProviderID           big.Int `json:"provider_id"`
	Price                int     `json:"price_per_byte"`
	Expiry               string  `json:"expiry_date"`
	QoS                  int     `json:"qos"`
	Signature            string  `json:"signature"`
	MerkleProof          string  `json:"merkle_proof"`
	FundedPaymentChannel bool    `json:"funded_payment_channel"` // Is this boolean?
}

