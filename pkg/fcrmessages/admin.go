package fcrmessages

import (
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// AdminGetReputationChallenge is the request from an admin client to a gateway to discover a client's reputation
type AdminGetReputationChallenge struct {
	ClientID nodeid.NodeID `json:"client_id"`
}

// EncodeAdminGetReputationChallenge is used to get the FCRMessage of AdminGetReputationChallenge
func EncodeAdminGetReputationChallenge(clientID *nodeid.NodeID) (*FCRMessage, error) {
	body, err := json.Marshal(AdminGetReputationChallenge{
		ClientID: *clientID,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       AdminGetReputationChallengeType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeAdminGetReputationChallenge is used to get the fields from FCRMessage of AdminGetReputationChallenge
func DecodeAdminGetReputationChallenge(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // client id
	error, // error
) {
	if fcrMsg.MessageType != AdminGetReputationChallengeType {
		return nil, fmt.Errorf("Message type mismatch")
	}
	msg := AdminGetReputationChallenge{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, err
	}
	return &msg.ClientID, nil
}

// AdminGetReputationResponse is the response to AdminGetReputationChallenge
type AdminGetReputationResponse struct {
	ClientID   nodeid.NodeID `json:"clientid"`
	Reputation int64         `json:"reputation"`
	Exists     bool          `json:"exists"`
}

// EncodeAdminGetReputationResponse is used to get the FCRMessage of AdminGetReputationResponse
func EncodeAdminGetReputationResponse(
	clientID *nodeid.NodeID,
	reputation int64,
	exists bool,
) (*FCRMessage, error) {
	body, err := json.Marshal(AdminGetReputationResponse{
		ClientID:   *clientID,
		Reputation: reputation,
		Exists:     exists,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       AdminGetReputationResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeAdminGetReputationResponse is used to get the fields from FCRMessage of AdminGetReputationResponse
func DecodeAdminGetReputationResponse(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // client id
	int64, // reputation
	bool, // exists
	error, // error
) {
	if fcrMsg.MessageType != AdminGetReputationResponseType {
		return nil, 0, false, fmt.Errorf("Message type mismatch")
	}
	msg := AdminGetReputationResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, false, err
	}
	return &msg.ClientID, msg.Reputation, msg.Exists, nil
}

// AdminSetReputationChallenge is the request from an admin client to a gateway to set a client's reputation
type AdminSetReputationChallenge struct {
	ClientID   nodeid.NodeID `json:"clientid"`
	Reputation int64         `json:"reputation"`
}

// EncodeAdminSetReputationChallenge is used to get the FCRMessage of AdminSetReputationChallenge
func EncodeAdminSetReputationChallenge(
	clientID *nodeid.NodeID,
	reputation int64,
) (*FCRMessage, error) {
	body, err := json.Marshal(AdminSetReputationChallenge{
		ClientID:   *clientID,
		Reputation: reputation,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       AdminSetReputationChallengeType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeAdminSetReputationChallenge is used to get the fields from FCRMessage of AdminSetReputationChallenge
func DecodeAdminSetReputationChallenge(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // client id
	int64, // reputation
	error, // error
) {
	if fcrMsg.MessageType != AdminSetReputationChallengeType {
		return nil, 0, fmt.Errorf("Message type mismatch")
	}
	msg := AdminSetReputationChallenge{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, err
	}
	return &msg.ClientID, msg.Reputation, nil
}

// AdminSetReputationResponse is the response to AdminSetReputationChallenge
type AdminSetReputationResponse struct {
	ClientID   nodeid.NodeID `json:"clientid"`
	Reputation int64         `json:"reputation"`
	Exists     bool          `json:"exists"`
}

// EncodeAdminSetReputationResponse is used to get the FCRMessage of AdminSetReputationResponse
func EncodeAdminSetReputationResponse(
	clientID *nodeid.NodeID,
	reputation int64,
	exists bool,
) (*FCRMessage, error) {
	body, err := json.Marshal(AdminSetReputationResponse{
		ClientID:   *clientID,
		Reputation: reputation,
		Exists:     exists,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       AdminSetReputationResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeAdminSetReputationResponse is used to get the fields from FCRMessage of AdminSetReputationResponse
func DecodeAdminSetReputationResponse(fcrMsg *FCRMessage) (
	*nodeid.NodeID, // client id
	int64, // reputation
	bool, // exists
	error, // error
) {
	if fcrMsg.MessageType != AdminSetReputationResponseType {
		return nil, 0, false, fmt.Errorf("Message type mismatch")
	}
	msg := AdminSetReputationResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return nil, 0, false, err
	}
	return &msg.ClientID, msg.Reputation, msg.Exists, nil
}

// AdminAcceptKeyChallenge is the request from an admin client to a gateway to generate an initial key pair.
type AdminAcceptKeyChallenge struct {
	PrivateKey        string `json:"privatekey"`
	PrivateKeyVersion uint32 `json:"privatekeyversion"`
}

// EncodeAdminAcceptKeyChallenge is used to get the FCRMessage of AdminAcceptKeysChallenge
func EncodeAdminAcceptKeyChallenge(
	string, // privatekey encoded as a hex string
	uint32, // privatekeyversion
) (*FCRMessage, error) {
	body, err := json.Marshal(AdminAcceptKeyChallenge{})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       AdminAcceptKeyChallengeType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeAdminAcceptKeyChallenge is used to get the fields from FCRMessage of AdminAcceptKeysChallenge
func DecodeAdminAcceptKeyChallenge(fcrMsg *FCRMessage) (string, uint32, error) {

	if fcrMsg.MessageType != AdminAcceptKeyChallengeType {
		return "", 0, fmt.Errorf("Message type mismatch")
	}
	msg := AdminAcceptKeyChallenge{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return "", 0, err
	}
	return msg.PrivateKey, msg.PrivateKeyVersion, nil
}

// AdminAcceptKeyResponse is the response to AdminAcceptKeysResponse
type AdminAcceptKeyResponse struct {
	Exists bool `json:"exists"`
}

// EncodeAdminAcceptKeyResponse is used to get the FCRMessage of AdminAcceptKeysResponse
// TODO: Set fields
func EncodeAdminAcceptKeyResponse(
	exists bool,
) (*FCRMessage, error) {
	body, err := json.Marshal(AdminAcceptKeyResponse{
		Exists: exists,
	})
	if err != nil {
		return nil, err
	}
	return &FCRMessage{
		MessageType:       AdminAcceptKeyResponseType,
		ProtocolVersion:   protocolVersion,
		ProtocolSupported: protocolSupported,
		MessageBody:       body,
	}, nil
}

// DecodeAdminAcceptKeyResponse is used to get the fields from FCRMessage of AdminAcceptKeysResponse
func DecodeAdminAcceptKeyResponse(fcrMsg *FCRMessage) (
	bool, // exists
	error, // error
) {
	if fcrMsg.MessageType != AdminAcceptKeyResponseType {
		return false, fmt.Errorf("Message type mismatch")
	}
	msg := AdminAcceptKeyResponse{}
	err := json.Unmarshal(fcrMsg.MessageBody, &msg)
	if err != nil {
		return false, err
	}
	return msg.Exists, nil
}
