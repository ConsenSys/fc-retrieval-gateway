package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ProviderPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer
type ProviderPublishGroupCIDRequest struct {
	Nonce      int64           `json:"nonce"`
	ProviderID nodeid.NodeID   `json:"provider_id"`
	Price      uint64          `json:"price_per_byte"`
	Expiry     int64           `json:"expiry_date"`
	QoS        uint64          `json:"qos"`
	PieceCIDs  []cid.ContentID `json:"piece_cids"`
}

// EncodeProviderPublishGroupCIDRequest is used to get the FCRMessage of ProviderPublishGroupCIDRequest
func EncodeProviderPublishGroupCIDRequest(
	nonce int64,
	providerID *nodeid.NodeID,
	price uint64,
	expiry int64,
	qos uint64,
	pieceCIDs []cid.ContentID,
) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderPublishGroupCIDRequest{
		Nonce:      nonce,
		ProviderID: *providerID,
		Price:      price,
		Expiry:     expiry,
		QoS:        qos,
		PieceCIDs:  pieceCIDs,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderPublishGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderPublishGroupCIDRequest is used to get the fields from FCRMessage of ProviderPublishGroupCIDRequest
func DecodeProviderPublishGroupCIDRequest(fcrMsg *FCRMessage) (
	int64, // nonce
	*nodeid.NodeID, // provider id
	uint64, // price
	int64, // expiry
	uint64, // qos
	[]cid.ContentID, // piece cids
	error, // error
) {
	if fcrMsg.MessageType != ProviderPublishGroupCIDRequestType {
		return 0, nil, 0, 0, 0, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderPublishGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return 0, nil, 0, 0, 0, nil, err
	}
	return msg.Nonce, &msg.ProviderID, msg.Price, msg.Expiry, msg.QoS, msg.PieceCIDs, nil
}

// ProviderDHTPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer using DHT
type ProviderDHTPublishGroupCIDRequest struct {
	Nonce           int64         `json:"nonce"`
	ProviderID      nodeid.NodeID `json:"provider_id"`
	NumOffers       int64         `json:"num_of_offers"`
	SingleCIDOffers []struct {
		Price     uint64        `json:"price_per_byte"`
		Expiry    int64         `json:"expiry_date"`
		QoS       uint64        `json:"qos"`
		Signature string        `json:"signature"`
		PieceCID  cid.ContentID `json:"piece_cid"`
	} `json:"single_cid_offers"`
}

// EncodeProviderDHTPublishGroupCIDRequest is used to get the FCRMessage of ProviderDHTPublishGroupCIDRequest
func EncodeProviderDHTPublishGroupCIDRequest(
	nonce int64,
	providerID *nodeid.NodeID,
	offers []cidoffer.CidGroupOffer, // TODO: Is this appropriate? Using group cid to represent single cid
) (*FCRMessage, error) {
	singleCIDOffers := make([]struct {
		Price     uint64        `json:"price_per_byte"`
		Expiry    int64         `json:"expiry_date"`
		QoS       uint64        `json:"qos"`
		Signature string        `json:"signature"`
		PieceCID  cid.ContentID `json:"piece_cid"`
	}, len(offers))
	for i, offer := range offers {
		singleCIDOffers[i].Price = offer.Price
		singleCIDOffers[i].Expiry = offer.Expiry
		singleCIDOffers[i].QoS = offer.QoS
		singleCIDOffers[i].Signature = offer.Signature
		singleCIDOffers[i].PieceCID = offer.Cids[0]
	}
	body, err := json.Marshal(ProviderDHTPublishGroupCIDRequest{
		Nonce:           nonce,
		ProviderID:      *providerID,
		NumOffers:       int64(len(offers)),
		SingleCIDOffers: singleCIDOffers,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderDHTPublishGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderDHTPublishGroupCIDRequest is used to get the fields from FCRMessage of ProviderDHTPublishGroupCIDRequest
func DecodeProviderDHTPublishGroupCIDRequest(fcrMsg *FCRMessage) (
	int64, // nonce
	*nodeid.NodeID, // provider id
	[]cidoffer.CidGroupOffer, // offers // TODO: Need to check if this is appropriate
	error, // error
) {
	if fcrMsg.MessageType != ProviderDHTPublishGroupCIDRequestType {
		return 0, nil, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderDHTPublishGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return 0, nil, nil, err
	}
	offers := make([]cidoffer.CidGroupOffer, msg.NumOffers)
	for i, singleCIDOffer := range msg.SingleCIDOffers {
		offers[i].NodeID = &msg.ProviderID
		offers[i].Cids = []cid.ContentID{singleCIDOffer.PieceCID}
		offers[i].Price = singleCIDOffer.Price
		offers[i].Expiry = singleCIDOffer.Expiry
		offers[i].QoS = singleCIDOffer.QoS
		offers[i].Signature = singleCIDOffer.Signature
		// TODO: What about the rest fields?
		// offers[i].MerkleProof
		// offers[i].FundedPaymentChannel
	}
	return msg.Nonce, &msg.ProviderID, offers, nil
}

// ProviderDHTPublishGroupCIDAck is the acknowledgement to ProviderDHTPublishGroupCIDRequest
type ProviderDHTPublishGroupCIDAck struct {
	Nonce     int64  `json:"nonce"`
	Signature string `json:"signature"`
}

// EncodeProviderDHTPublishGroupCIDAck is used to get the FCRMessage of ProviderDHTPublishGroupCIDAck
func EncodeProviderDHTPublishGroupCIDAck(
	nonce int64,
	signature string,
) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderDHTPublishGroupCIDAck{
		Nonce:     nonce,
		Signature: signature,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderDHTPublishGroupCIDAckType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderDHTPublishGroupCIDAck is used to get the fields from FCRMessage of ProviderDHTPublishGroupCIDAck
func DecodeProviderDHTPublishGroupCIDAck(fcrMsg *FCRMessage) (
	int64, // nonce
	string, // signature
	error, // error
) {
	if fcrMsg.MessageType != ProviderDHTPublishGroupCIDAckType {
		return 0, "", fmt.Errorf("Message type mismatch")
	}
	msg := ProviderDHTPublishGroupCIDAck{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return 0, "", err
	}
	return msg.Nonce, msg.Signature, nil
}
