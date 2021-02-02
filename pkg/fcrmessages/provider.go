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
	Signature  string          `json:"signature"`
}

// EncodeProviderPublishGroupCIDRequest is used to get the FCRMessage of ProviderPublishGroupCIDRequest
func EncodeProviderPublishGroupCIDRequest(
	nonce int64,
	offer *cidoffer.CidGroupOffer,
) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderPublishGroupCIDRequest{
		Nonce:      nonce,
		ProviderID: *offer.NodeID,
		Price:      offer.Price,
		Expiry:     offer.Expiry,
		QoS:        offer.QoS,
		PieceCIDs:  offer.Cids,
		Signature:  offer.Signature,
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

// ProviderDHTPublishGroupCIDRequest is the request from provider to gateway to publish group cid offer using DHT
type ProviderDHTPublishGroupCIDRequest struct {
	Nonce           int64         `json:"nonce"`
	ProviderID      nodeid.NodeID `json:"provider_id"`
	NumOffers       int64         `json:"num_of_offers"`
	SingleCIDOffers []struct {
		Price     uint64        `json:"price_per_byte"`
		Expiry    int64         `json:"expiry_date"`
		QoS       uint64        `json:"qos"`
		PieceCID  cid.ContentID `json:"piece_cid"`
		Signature string        `json:"signature"`
	} `json:"single_cid_offers"`
}

// EncodeProviderDHTPublishGroupCIDRequest is used to get the FCRMessage of ProviderDHTPublishGroupCIDRequest
func EncodeProviderDHTPublishGroupCIDRequest(
	nonce int64,
	providerID *nodeid.NodeID,
	offers []cidoffer.CidGroupOffer,
) (*FCRMessage, error) {
	singleCIDOffers := make([]struct {
		Price     uint64        `json:"price_per_byte"`
		Expiry    int64         `json:"expiry_date"`
		QoS       uint64        `json:"qos"`
		PieceCID  cid.ContentID `json:"piece_cid"`
		Signature string        `json:"signature"`
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
	[]cidoffer.CidGroupOffer, // offers
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
		offer, err := cidoffer.NewCidGroupOffer(
			&msg.ProviderID,
			&[]cid.ContentID{singleCIDOffer.PieceCID},
			singleCIDOffer.Price,
			singleCIDOffer.Expiry,
			singleCIDOffer.QoS)
		if err != nil {
			return 0, nil, nil, err
		}
		offers[i].Signature = singleCIDOffer.Signature
		offers = append(offers, *offer)
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
