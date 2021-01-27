package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ClientEstablishmentRequest is the request from client to gateway to establish connection
type ClientEstablishmentRequest struct {
	ClientID  nodeid.NodeID `json:"client_id"`
	Challenge string        `json:"challenge"`
	TTL       int64         `json:"ttl"`
}

// EncodeClientEstablishmentRequest is used to get the FCRMessage of ClientEstablishmentRequest
func EncodeClientEstablishmentRequest(
	clientID *nodeid.NodeID,
	challenge string,
	ttl int64,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientEstablishmentRequest{
		ClientID:  *clientID,
		Challenge: challenge,
		TTL:       ttl,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientEstablishmentRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientEstablishmentRequest is used to get the fields from FCRMessage of ClientEstablishmentRequest
func DecodeClientEstablishmentRequest(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // client id
	string, // challenge
	int64, // ttl
	error, // error
) {
	if fcrMsg.MessageType != ClientEstablishmentRequestType {
		return nil, "", 0, fmt.Errorf("Message type mismatch")
	}
	msg := ClientEstablishmentRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, "", 0, err
	}
	return &msg.ClientID, msg.Challenge, msg.TTL, nil
}

// ClientEstablishmentResponse is the response to ClientEstablishmentRequest
type ClientEstablishmentResponse struct {
	GatewayID nodeid.NodeID `json:"gateway_id"`
	Challenge string        `json:"challenge"`
	Signature string        `json:"signature"`
}

// EncodeClientEstablishmentResponse is used to get the FCRMessage of ClientEstablishmentResponse
func EncodeClientEstablishmentResponse(
	gatewayID *nodeid.NodeID,
	challenge string,
	signature string,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientEstablishmentResponse{
		GatewayID: *gatewayID,
		Challenge: challenge,
		Signature: signature,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientEstablishmentResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientEstablishmentResponse is used to get the fields from FCRMessage of ClientEstablishmentResponse
func DecodeClientEstablishmentResponse(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // gateway id
	string, // challenge
	string, // signature
	error, // error
) {
	if fcrMsg.MessageType != ClientEstablishmentResponseType {
		return nil, "", "", fmt.Errorf("Message type mismatch")
	}
	msg := ClientEstablishmentResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, "", "", err
	}
	return &msg.GatewayID, msg.Challenge, msg.Signature, nil
}

// ClientStandardDiscoverRequest is the requset from client to gateway to ask for cid offer
type ClientStandardDiscoverRequest struct {
	PieceCID cid.ContentID `json:"piece_cid"`
	Nonce    int64         `json:"nonce"`
	TTL      int64         `json:"ttl"`
}

// EncodeClientStandardDiscoverRequest is used to get the FCRMessage of ClientStandardDiscoverRequest
func EncodeClientStandardDiscoverRequest(
	pieceCID *cid.ContentID,
	nonce int64,
	ttl int64,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientStandardDiscoverRequest{
		PieceCID: *pieceCID,
		Nonce:    nonce,
		TTL:      ttl,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientStandardDiscoverRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientStandardDiscoverRequest is used to get the fields from FCRMessage of ClientStandardDiscoverRequest
func DecodeClientStandardDiscoverRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	int64, // ttl
	error, // error
) {
	if fcrMsg.MessageType != ClientStandardDiscoverRequestType {
		return nil, 0, 0, fmt.Errorf("Message type mismatch")
	}
	msg := ClientStandardDiscoverRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, 0, err
	}
	return &msg.PieceCID, msg.Nonce, msg.TTL, nil
}

// ClientStandardDiscoverResponse is the response to ClientStandardDiscoverResponse
type ClientStandardDiscoverResponse struct {
	PieceCID     cid.ContentID         `json:"piece_cid"`
	Nonce        int64                 `json:"nonce"`
	Found        bool                  `json:"found"`
	CIDGroupInfo []CIDGroupInformation `json:"cid_group_information"`
}

// EncodeClientStandardDiscoverResponse is used to get the FCRMessage of ClientStandardDiscoverResponse
func EncodeClientStandardDiscoverResponse(
	pieceCID *cid.ContentID,
	nonce int64,
	found bool,
	offers []*cidoffer.CidGroupOffer,
) (*FCRMessage, error) {
	cidGroupInfo := make([]CIDGroupInformation, 0)
	if found {
		for _, offer := range offers {
			cidGroupInfo = append(cidGroupInfo, CIDGroupInformation{
				ProviderID:           *offer.NodeID,
				Price:                offer.Price,
				Expiry:               offer.Expiry,
				QoS:                  offer.QoS,
				Signature:            offer.Signature,
				MerkleProof:          offer.MerkleProof,
				FundedPaymentChannel: offer.FundedPaymentChannel,
			})
		}
	}
	body, err := json.Marshal(ClientStandardDiscoverResponse{
		PieceCID:     *pieceCID,
		Nonce:        nonce,
		Found:        found,
		CIDGroupInfo: cidGroupInfo,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientStandardDiscoverResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientStandardDiscoverResponse is used to get the fields from FCRMessage of ClientStandardDiscoverResponse
func DecodeClientStandardDiscoverResponse(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	bool, // found
	[]*cidoffer.CidGroupOffer, // offers
	error, // error
) {
	if fcrMsg.MessageType != ClientStandardDiscoverResponseType {
		return nil, 0, false, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ClientStandardDiscoverResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, false, nil, err
	}
	offers := make([]*cidoffer.CidGroupOffer, 0)
	if msg.Found {
		for _, offerInfo := range msg.CIDGroupInfo {
			offers = append(offers, &cidoffer.CidGroupOffer{
				NodeID:               &offerInfo.ProviderID,
				Cids:                 []cid.ContentID{msg.PieceCID},
				Price:                offerInfo.Price,
				Expiry:               offerInfo.Expiry,
				QoS:                  offerInfo.QoS,
				Signature:            offerInfo.Signature,
				MerkleProof:          offerInfo.MerkleProof,
				FundedPaymentChannel: offerInfo.FundedPaymentChannel,
			})
		}
	}
	return &msg.PieceCID, msg.Nonce, msg.Found, offers, nil
}

// ClientDHTDiscoverRequest is the request from client to gateway to ask for cid offer using DHT
type ClientDHTDiscoverRequest struct {
	PieceCID           cid.ContentID `json:"piece_cid"`
	Nonce              int64         `json:"nonce"`
	TTL                int64         `json:"ttl"`
	NumDHT             int64         `json:"num_dht"`
	IncrementalResults bool          `json:"incremental_results"`
}

// EncodeClientDHTDiscoverRequest is used to get the FCRMessage of ClientDHTDiscoverRequest
func EncodeClientDHTDiscoverRequest(
	pieceCID *cid.ContentID,
	nonce int64,
	ttl int64,
	numDHT int64,
	incrementalResults bool,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientDHTDiscoverRequest{
		PieceCID:           *pieceCID,
		Nonce:              nonce,
		TTL:                ttl,
		NumDHT:             numDHT,
		IncrementalResults: incrementalResults,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientDHTDiscoverRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientDHTDiscoverRequest is used to get the fields from FCRMessage of ClientDHTDiscoverRequest
func DecodeClientDHTDiscoverRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	int64, // nonce
	int64, // ttl
	int64, // num dht
	bool, // incremental results
	error, // error
) {
	if fcrMsg.MessageType != ClientDHTDiscoverRequestType {
		return nil, 0, 0, 0, false, fmt.Errorf("Message type mismatch")
	}
	msg := ClientDHTDiscoverRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, 0, 0, false, err
	}
	return &msg.PieceCID, msg.Nonce, msg.TTL, msg.NumDHT, msg.IncrementalResults, nil
}

// ClientDHTDiscoverResponse is the response to ClientDHTDiscoverRequest
type ClientDHTDiscoverResponse struct {
	Contacted     []FCRMessage    `json:"contacted_gateways"`
	UnContactable []nodeid.NodeID `json:"uncontactable_gateways"`
	Nonce         int64           `json:"nonce"`
}

// EncodeClientDHTDiscoverResponse is used to get the FCRMessage of ClientDHTDiscoverResponse
func EncodeClientDHTDiscoverResponse(
	contacted []FCRMessage,
	unContactable []nodeid.NodeID,
	nonce int64,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientDHTDiscoverResponse{
		Contacted:     contacted,
		UnContactable: unContactable,
		Nonce:         nonce,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientDHTDiscoverResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientDHTDiscoverResponse is used to get the fields from FCRMessage of ClientDHTDiscoverResponse
func DecodeClientDHTDiscoverResponse(fcrMsg *FCRMessage) (
	[]FCRMessage, // contacted
	[]nodeid.NodeID, // uncontactable
	int64, // nonce
	error, // error
) {
	if fcrMsg.MessageType != ClientDHTDiscoverResponseType {
		return nil, nil, 0, fmt.Errorf("Message type mismatch")
	}
	msg := ClientDHTDiscoverResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, nil, 0, err
	}
	return msg.Contacted, msg.UnContactable, msg.Nonce, nil
}

// ClientCIDGroupPublishDHTAckRequest is the request from client to provider to request the signed ack of a cid group publish
type ClientCIDGroupPublishDHTAckRequest struct {
	PieceCID  cid.ContentID `json:"piece_cid"`
	GatewayID nodeid.NodeID `json:"gateway_id"`
}

// EncodeClientCIDGroupPublishDHTAckRequest is used to get the FCRMessage of ClientCIDGroupPublishDHTAckRequest
func EncodeClientCIDGroupPublishDHTAckRequest(
	pieceCID *cid.ContentID,
	gatewayID *nodeid.NodeID,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientCIDGroupPublishDHTAckRequest{
		PieceCID:  *pieceCID,
		GatewayID: *gatewayID,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientCIDGroupPublishDHTAckRequestType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientCIDGroupPublishDHTAckRequest is used to get the fields from FCRMessage of ClientCIDGroupPublishDHTAckRequest
func DecodeClientCIDGroupPublishDHTAckRequest(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	*nodeid.NodeID, // gateway id
	error, // error
) {
	if fcrMsg.MessageType != ClientCIDGroupPublishDHTAckRequestType {
		return nil, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ClientCIDGroupPublishDHTAckRequest{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, nil, err
	}
	return &msg.PieceCID, &msg.GatewayID, nil
}

// ClientCIDGroupPublishDHTAckResponse is the response to ClientCIDGroupPublishDHTAckRequest
type ClientCIDGroupPublishDHTAckResponse struct {
	PieceCID                cid.ContentID `json:"piece_cid"`
	GatewayID               nodeid.NodeID `json:"gateway_id"`
	Found                   bool          `json:"found"`
	CIDGroupPublishToDHT    FCRMessage    `json:"cid_group_publish_to_dht"`
	CIDGroupPublishToDHTAck FCRMessage    `json:"cid_group_publish_to_dht_ack"`
}

// EncodeClientCIDGroupPublishDHTAckResponse is used to get the FCRMessage of ClientCIDGroupPublishDHTAckResponse
func EncodeClientCIDGroupPublishDHTAckResponse(
	pieceCID *cid.ContentID,
	gatewayID *nodeid.NodeID,
	found bool,
	cidGroupPublishToDHT *FCRMessage,
	cidGroupPublishToDHTAck *FCRMessage,
) (*FCRMessage, error) {
	body, err := json.Marshal(ClientCIDGroupPublishDHTAckResponse{
		PieceCID:                *pieceCID,
		GatewayID:               *gatewayID,
		Found:                   found,
		CIDGroupPublishToDHT:    *cidGroupPublishToDHT,
		CIDGroupPublishToDHTAck: *cidGroupPublishToDHTAck,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       ClientCIDGroupPublishDHTAckResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeClientCIDGroupPublishDHTAckResponse is used to get the fields from FCRMessage of ClientCIDGroupPublishDHTAckResponse
func DecodeClientCIDGroupPublishDHTAckResponse(fcrMsg *FCRMessage) (
	*cid.ContentID, // piece cid
	*nodeid.NodeID, // gateway id
	bool, // found
	*FCRMessage, // cid group publish to dht
	*FCRMessage, // cid group publish to dht ack
	error, // error
) {
	if fcrMsg.MessageType != ClientCIDGroupPublishDHTAckResponseType {
		return nil, nil, false, nil, nil, fmt.Errorf("Message type mismatch")
	}
	msg := ClientCIDGroupPublishDHTAckResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, nil, false, nil, nil, err
	}
	return &msg.PieceCID, &msg.GatewayID, msg.Found, &msg.CIDGroupPublishToDHT, &msg.CIDGroupPublishToDHTAck, nil
}
