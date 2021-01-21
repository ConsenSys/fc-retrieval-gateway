package fcrmessages

import (
	"encoding/json"
)

const (
	defaultProtocolVersion            = 1
	defaultAlternativeProtocolVersion = 1
)

var protocolVersion int32 = defaultProtocolVersion
var protocolSupported []int32 = []int32{defaultProtocolVersion, defaultAlternativeProtocolVersion}

// FCRMessage is the message used in communication between filecoin retrieval entities
type FCRMessage struct {
	MessageType       int32   `json:"message_type"`
	ProtocolVersion   int32   `json:"protocol_version"`
	ProtocolSupported []int32 `json:"protocol_supported"`
	MessageBody       []byte  `json:"message_body"`
	Signature         string  `json:"message_signature"`
}

// FCRMsgToBytes converts a FCRMessage to bytes
func FCRMsgToBytes(fcrMsg *FCRMessage) ([]byte, error) {
	return json.Marshal(fcrMsg)
}

// FCRMsgFromBytes converts a bytes to FCRMessage
func FCRMsgFromBytes(data []byte) (*FCRMessage, error) {
	res := FCRMessage{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetProtocolVersion gets the current protocol version of all messages
func GetProtocolVersion() (int32, []int32) {
	return protocolVersion, protocolSupported
}

// SetProtocolVersion sets the current protocol version of all messages
func SetProtocolVersion(newProtocolVersion int32, newProtocolSupported []int32) {
	protocolVersion = newProtocolVersion
	protocolSupported = newProtocolSupported
}

// EncodeXX is used to get the FCRMessage of XX
// DecodeXX is used to get the fields from FCRMessage of XX
