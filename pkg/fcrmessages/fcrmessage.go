package fcrmessages

import (
	"crypto/sha512"
	"encoding/json"
)

const (
	defaultProtocolVersion            = 1
	defaultAlternativeProtocolVersion = 1
	CidGroupOfferDigestSize						= sha512.Size256
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

// GetMessageType is used to get the message type of the message
func (fcrMsg *FCRMessage) GetMessageType() int32 {
	return fcrMsg.MessageType
}

// GetProtocolVersion is used to get the protocol version of peer
func (fcrMsg *FCRMessage) GetProtocolVersion() int32 {
	return fcrMsg.ProtocolVersion
}

// GetProtocolSupported is used to get the protocol supported of peer
func (fcrMsg *FCRMessage) GetProtocolSupported() []int32 {
	return fcrMsg.GetProtocolSupported()
}

// GetMessageBody is used to get the message body
func (fcrMsg *FCRMessage) GetMessageBody() []byte {
	return fcrMsg.MessageBody
}

// VerifySignature is used to verify the signature
func (fcrMsg *FCRMessage) VerifySignature(verify func(sig string, msg interface{}) (bool, error)) (bool, error) {
	// Clear signature
	sig := fcrMsg.Signature
	fcrMsg.Signature = ""
	res, err := verify(sig, fcrMsg)
	if err != nil {
		return false, err
	}
	// Recover signature
	fcrMsg.Signature = sig
	return res, nil
}

// SignMessage is used to sign the message
func (fcrMsg *FCRMessage) SignMessage(sign func(msg interface{}) (string, error)) error {
	sig, err := sign(fcrMsg)
	if err != nil {
		return err
	}
	fcrMsg.Signature = sig
	return nil
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

// DumpMessage display a message with ASCII and hex
func (fcrMsg *FCRMessage) DumpMessage() string { 
	msgBytes, err := FCRMsgToBytes(fcrMsg)
	if err != nil {
		return "Error processing message"
	}
	return dumpMessage(msgBytes)
}

