package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ProviderPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer
type ProviderRegistrationRequest struct {
	Nonce      int64           `json:"nonce"`
	ProviderID nodeid.NodeID   `json:"provider_id"`
	Address    	string          `json:"price_per_byte"`
	NetworkInfo	string           `json:"expiry_date"`
	RegionCode	string          `json:"qos"`
	PieceCIDs  []cid.ContentID `json:"piece_cids"`
	Signature  string          `json:"signature"`
}

// EncodeProviderRegistrationRequest is used to get the FCRMessage of ProviderRegistrationRequest
func EncodeProviderRegistrationRequest() (*FCRMessage, error) {
	return &FCRMessage{
		MessageType:       ProviderAdminRegistrationRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       nil,
	}, nil
}

// DecodeProviderRegistrationRequest is used to get the fields from FCRMessage of ProviderRegistrationRequest
func DecodeProviderRegistrationRequest(fcrMsg *FCRMessage) (
	int64, // nonce
	*cidoffer.CidGroupOffer, // offer
	error, // error
) {
	if fcrMsg.MessageType != ProviderPublishGroupCIDRequestType {
		return 0, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderPublishGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return 0, nil, err
	}
	offer, err := cidoffer.NewCidGroupOffer(
		&msg.ProviderID,
		&msg.PieceCIDs,
		msg.Price,
		msg.Expiry,
		msg.QoS)
	if err != nil {
		return 0, nil, err
	}
	// Set signature
	offer.Signature = msg.Signature
	return msg.Nonce, offer, nil
}
