package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmerkletrie"
)

// GatewaySingleCIDOfferPublishRequest is the request from gateway to provider during start-up asking for cid offers
type GatewaySingleCIDOfferPublishRequest struct {
	CIDMin             cid.ContentID                `json:"cid_min"`
	CIDMax             cid.ContentID                `json:"cid_max"`
	BlockHash          string                       `json:"block_hash"`
	TransactionReceipt string                       `json:"transaction_receipt"`
	MerkleRoot         string                       `json:"merkle_root"`
	MerkleProof        fcrmerkletrie.FCRMerkleProof `json:"merkle_proof"`
}

// EncodeGatewaySingleCIDOfferPublishRequest is used to get the FCRMessage of GatewaySingleCIDOfferPublishRequest
func EncodeGatewaySingleCIDOfferPublishRequest(
	cidMin *cid.ContentID,
	cidMax *cid.ContentID,
	blockHash string,
	transactionReceipt string,
	merkleRoot string,
	merkleProof *fcrmerkletrie.FCRMerkleProof,
) (*FCRMessage, error) {
	body, err := json.Marshal(GatewaySingleCIDOfferPublishRequest{
		CIDMin:             *cidMin,
		CIDMax:             *cidMax,
		BlockHash:          blockHash,
		TransactionReceipt: transactionReceipt,
		MerkleRoot:         merkleRoot,
		MerkleProof:        *merkleProof,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       GatewaySingleCIDOfferPublishRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeGatewaySingleCIDOfferPublishRequest is used to get the fields from FCRMessage of GatewaySingleCIDOfferPublishRequest
func DecodeGatewaySingleCIDOfferPublishRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // cid min
	*cid.ContentID, // cid max
	string, // block hash
	string, // transaction receipt
	string, // merkle root
	*fcrmerkletrie.FCRMerkleProof, // merkle proof
	error, // error
) {
	if fcrMsg.MessageType != GatewaySingleCIDOfferPublishRequestType {
		return nil, nil, "", "", "", nil, fmt.Errorf("Message type mismatch")
	}
	msg := GatewaySingleCIDOfferPublishRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, nil, "", "", "", nil, err
	}
	return &msg.CIDMin, &msg.CIDMax, msg.BlockHash, msg.TransactionReceipt, msg.MerkleRoot, &msg.MerkleProof, nil
}

// GatewaySingleCIDOfferPublishResponse is the repsonse to GatewaySingleCIDOfferPublishRequest
type GatewaySingleCIDOfferPublishResponse struct {
	PublishedGroupCIDs []FCRMessage `json:"published_dht_cid_offers"`
}

// EncodeGatewaySingleCIDOfferPublishResponse is used to get the FCRMessage of GatewaySingleCIDOfferPublishResponse
func EncodeGatewaySingleCIDOfferPublishResponse(
	publishedGroupCIDs []FCRMessage,
) (*FCRMessage, error) {
	body, err := json.Marshal(GatewaySingleCIDOfferPublishResponse{
		PublishedGroupCIDs: publishedGroupCIDs,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       GatewaySingleCIDOfferPublishResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeGatewaySingleCIDOfferPublishResponse is used to get the fields from FCRMessage of GatewaySingleCIDOfferPublishResponse
func DecodeGatewaySingleCIDOfferPublishResponse(fcrMsg *FCRMessage) (
	[]FCRMessage, // published dht cid offers
	error, // error
) {
	if fcrMsg.MessageType != GatewaySingleCIDOfferPublishResponseType {
		return nil, fmt.Errorf("Message type mismatch")
	}
	msg := GatewaySingleCIDOfferPublishResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, err
	}
	return msg.PublishedGroupCIDs, nil
}

// GatewaySingleCIDOfferPublishResponseAck is the acknowledgement to GatewaySingleCIDOfferPublishResponse
type GatewaySingleCIDOfferPublishResponseAck struct {
	CIDOffersAck []FCRMessage `json:"published_dht_cid_offers_ack"`
}

// EncodeGatewaySingleCIDOfferPublishResponseAck is used to get the FCRMessage of GatewaySingleCIDOfferPublishResponseAck
func EncodeGatewaySingleCIDOfferPublishResponseAck(
	cidOffersAck []FCRMessage,
) (*FCRMessage, error) {
	body, err := json.Marshal(GatewaySingleCIDOfferPublishResponseAck{
		CIDOffersAck: cidOffersAck,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       GatewaySingleCIDOfferPublishResponseAckType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeGatewaySingleCIDOfferPublishResponseAck is used to get the fields from FCRMessage of GatewaySingleCIDOfferPublishResponseAck
func DecodeGatewaySingleCIDOfferPublishResponseAck(fcrMsg *FCRMessage) (
	[]FCRMessage, // published dht cid offers ack
	error, // error
) {
	if fcrMsg.MessageType != GatewaySingleCIDOfferPublishResponseAckType {
		return nil, fmt.Errorf("Message type mismatch")
	}
	msg := GatewaySingleCIDOfferPublishResponseAck{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, err
	}
	return msg.CIDOffersAck, nil
}

// GatewayDHTDiscoverRequest is the request from gateway to gateway to discover cid offer
type GatewayDHTDiscoverRequest struct {
	PieceCID cid.ContentID `json:"piece_cid"`
	Nonce    int64         `json:"nonce"`
	TTL      int64         `json:"ttl"`
}

// EncodeGatewayDHTDiscoverRequest is used to get the FCRMessage of GatewayDHTDiscoverRequest
func EncodeGatewayDHTDiscoverRequest(
	pieceCID *cid.ContentID,
	nonce int64,
	ttl int64,
) (*FCRMessage, error) {
	body, err := json.Marshal(GatewayDHTDiscoverRequest{
		PieceCID: *pieceCID,
		Nonce:    nonce,
		TTL:      ttl,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       GatewayDHTDiscoverRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeGatewayDHTDiscoverRequest is used to get the fields from FCRMessage of GatewayDHTDiscoverRequest
func DecodeGatewayDHTDiscoverRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	int64, // ttl
	error, // error
) {
	if fcrMsg.MessageType != GatewayDHTDiscoverRequestType {
		return nil, 0, 0, fmt.Errorf("Message type mismatch")
	}
	msg := GatewayDHTDiscoverRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, 0, err
	}
	return &msg.PieceCID, msg.Nonce, msg.TTL, nil
}

// GatewayDHTDiscoverResponse is the response to GatewayDHTDiscoverRequest
type GatewayDHTDiscoverResponse struct {
	PieceCID     cid.ContentID         `json:"piece_cid"`
	Nonce        int64                 `json:"nonce"`
	Found        bool                  `json:"found"`
	CIDGroupInfo []CIDGroupInformation `json:"cid_group_information"`
}

// EncodeGatewayDHTDiscoverResponse is used to get the FCRMessage of GatewayDHTDiscoverResponse
func EncodeGatewayDHTDiscoverResponse(
	pieceCID *cid.ContentID,
	nonce int64,
	found bool,
	offers []*cidoffer.CidGroupOffer,
	roots []string,
	proofs []fcrmerkletrie.FCRMerkleProof,
	fundedPaymentChannel []bool,
) (*FCRMessage, error) {
	cidGroupInfo := make([]CIDGroupInformation, len(offers))
	if found {
		for i := 0; i < len(offers); i++ {
			offer := offers[i]
			cidGroupInfo[i] = CIDGroupInformation{
				ProviderID:           *offer.NodeID,
				Price:                offer.Price,
				Expiry:               offer.Expiry,
				QoS:                  offer.QoS,
				Signature:            offer.Signature,
				MerkleRoot:           roots[i],
				MerkleProof:          proofs[i],
				FundedPaymentChannel: fundedPaymentChannel[i],
			}
		}
	}
	body, err := json.Marshal(GatewayDHTDiscoverResponse{
		PieceCID:     *pieceCID,
		Nonce:        nonce,
		Found:        found,
		CIDGroupInfo: cidGroupInfo,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       GatewayDHTDiscoverResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeGatewayDHTDiscoverResponse is used to get the fields from FCRMessage of GatewayDHTDiscoverResponse
func DecodeGatewayDHTDiscoverResponse(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	bool, // found
	[]cidoffer.CidGroupOffer, // offers
	[]string, // merkle roots
	[]fcrmerkletrie.FCRMerkleProof, // merkle proofs
	[]bool, // funded payment channel
	error, // error
) {
	if fcrMsg.MessageType != GatewayDHTDiscoverResponseType {
		return nil, 0, false, nil, nil, nil, nil, fmt.Errorf("Message type mimatch")
	}
	msg := GatewayDHTDiscoverResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, false, nil, nil, nil, nil, err
	}
	offers := make([]cidoffer.CidGroupOffer, 0)
	roots := make([]string, 0)
	proofs := make([]fcrmerkletrie.FCRMerkleProof, 0)
	fundedPaymentChannel := make([]bool, 0)
	if msg.Found {
		for _, offerInfo := range msg.CIDGroupInfo {
			offer, err := cidoffer.NewCidGroupOffer(
				&offerInfo.ProviderID,
				&[]cid.ContentID{msg.PieceCID},
				offerInfo.Price,
				offerInfo.Expiry,
				offerInfo.QoS)
			if err != nil {
				return nil, 0, false, nil, nil, nil, nil, err
			}
			// Set signature
			offer.Signature = offerInfo.Signature
			offers = append(offers, *offer)
			roots = append(roots, offerInfo.MerkleRoot)
			proofs = append(proofs, offerInfo.MerkleProof)
			fundedPaymentChannel = append(fundedPaymentChannel, offerInfo.FundedPaymentChannel)
		}
	}
	return &msg.PieceCID, msg.Nonce, msg.Found, offers, roots, proofs, fundedPaymentChannel, nil
}
