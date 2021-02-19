package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ProviderAdminGetGroupCIDRequest is the requset from client to gateway to ask for cid offer
type ProviderAdminGetGroupCIDRequest struct {
	GatewayIDs	[]nodeid.NodeID `json:"gateway_id"`
}

// EncodeProviderAdminGetGroupCIDRequest is used to get the FCRMessage of ProviderAdminGetGroupCIDRequest
func EncodeProviderAdminGetGroupCIDRequest(
	gatewayIDs []nodeid.NodeID,
) (*FCRMessage, error) {
	fmt.Printf("EncodeProviderAdminGetGroupCIDRequest\n")
	body, err := json.Marshal(ProviderAdminGetGroupCIDRequest{
		GatewayIDs: gatewayIDs,
	})
	fmt.Printf("Body: %+v\n", body)
	if err != nil {
		fmt.Printf("Error when marshalling request: %+v\n", err)
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminGetGroupCIDRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminGetGroupCIDRequest is used to get the fields from FCRMessage of ProviderAdminGetGroupCIDRequest
func DecodeProviderAdminGetGroupCIDRequest(fcrMsg *FCRMessage) (
	[]nodeid.NodeID, // piece cids
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminGetGroupCIDRequestType {
		return nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminGetGroupCIDRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, err
	}
	return msg.GatewayIDs, nil
}

// ProviderAdminGetGroupCIDResponse is the response to ProviderAdminGetGroupCIDResponse
type ProviderAdminGetGroupCIDResponse struct {
	Found        	bool                  `json:"found"`
	CIDGroupInfo 	[]CIDGroupInformation `json:"cid_group_information"`
}

// EncodeProviderAdminGetGroupCIDResponse is used to get the FCRMessage of ProviderAdminGetGroupCIDResponse
func EncodeProviderAdminGetGroupCIDResponse(
	found bool,
	offers []*cidoffer.CidGroupOffer,
	roots []string,
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
				FundedPaymentChannel: fundedPaymentChannel[i],
			}
		}
	}
	body, err := json.Marshal(ProviderAdminGetGroupCIDResponse{
		Found:        found,
		CIDGroupInfo: cidGroupInfo,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ProviderAdminGetGroupCIDResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeProviderAdminGetGroupCIDResponse is used to get the fields from FCRMessage of ProviderAdminGetGroupCIDResponse
func DecodeProviderAdminGetGroupCIDResponse(fcrMsg *FCRMessage) (
	bool, // found
	[]CIDGroupInformation, // offers
	error, // error
) {
	if fcrMsg.MessageType != ProviderAdminGetGroupCIDResponseType {
		return false, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ProviderAdminGetGroupCIDResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return false, nil, err
	}
	return msg.Found, msg.CIDGroupInfo, nil
}