package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmerkletree"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ProviderGetGroupCIDRequest is the requset from client to gateway to ask for cid offer
type ProviderGetGroupCIDRequest struct {
	ProviderID	nodeid.NodeID `json:"provider_id"`
}

// EncodeProviderGetGroupCIDRequest is used to get the FCRMessage of ProviderGetGroupCIDRequest
func EncodeProviderGetGroupCIDRequest(
	providerID *nodeid.NodeID,
) (*FCRMessage, error) {
	body, err := json.Marshal(ProviderGetGroupCIDRequest{
		ProviderID: *providerID,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderGetGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderGetGroupCIDRequest is used to get the fields from FCRMessage of ProviderGetGroupCIDRequest
func DecodeProviderGetGroupCIDRequest(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // piece cid
	error, // error
) {
	if fcrMsg.MessageType != ProviderGetGroupCIDRequestType {
		return nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderGetGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, err
	}
	return &msg.ProviderID, nil
}

// ProviderGetGroupCIDResponse is the response to ProviderGetGroupCIDResponse
type ProviderGetGroupCIDResponse struct {
	ProviderID		nodeid.NodeID 				`json:"provider_id"`
	Found        	bool                  `json:"found"`
	CIDGroupInfo 	[]CIDGroupInformation `json:"cid_group_information"`
}

// EncodeProviderGetGroupCIDResponse is used to get the FCRMessage of ProviderGetGroupCIDResponse
func EncodeProviderGetGroupCIDResponse(
	providerID *nodeid.NodeID,
	found bool,
	offers []*cidoffer.CidGroupOffer,
	roots []string,
	proofs []fcrmerkletree.FCRMerkleProof,
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
	body, err := json.Marshal(ProviderGetGroupCIDResponse{
		ProviderID:   *providerID,
		Found:        found,
		CIDGroupInfo: cidGroupInfo,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderGetGroupCIDResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderGetGroupCIDResponse is used to get the fields from FCRMessage of ProviderGetGroupCIDResponse
func DecodeProviderGetGroupCIDResponse(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // provider id
	bool, // found
	[]cidoffer.CidGroupOffer, // offers
	error, // error
) {
	if fcrMsg.MessageType != ProviderGetGroupCIDResponseType {
		return nil, false, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderGetGroupCIDResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, false, nil, err
	}
	offers := make([]cidoffer.CidGroupOffer, 0)
	return &msg.ProviderID, msg.Found, offers, nil
}