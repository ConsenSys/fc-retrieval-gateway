package tcpcomms

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
)

// IsTimeoutError checks if the given error is a timeout error
func IsTimeoutError(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr.Timeout()
}

// ReadTCPMessage read the tcp message from a given connection
func ReadTCPMessage(conn net.Conn, timeout time.Duration) (byte, []byte, error) {
	// Initialise a reader
	reader := bufio.NewReader(conn)
	// Read the length
	length := make([]byte, 4)
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	_, err := io.ReadFull(reader, length)
	if err != nil {
		return 0, nil, err
	}
	// Read the request
	request := make([]byte, int(binary.BigEndian.Uint32(length)))
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	_, err = io.ReadFull(reader, request)
	if err != nil {
		return 0, nil, err
	}
	return request[0], request[1:], nil
}

// SendTCPMessage sends a tcp message to a given connection
func SendTCPMessage(conn net.Conn, msgType byte, data []byte, timeout time.Duration) error {
	// Initialise a writer
	writer := bufio.NewWriter(conn)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(1+len(data)))
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	_, err := writer.Write(append(append(length, msgType), data...))
	if err != nil {
		return err
	}
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	return writer.Flush()
}

// SendProtocolMismatch sends a protocol mistmatch message to a given connection
func SendProtocolMismatch(conn net.Conn, timeout time.Duration) error {
	data, _ := json.Marshal(messages.ProtocolMismatchResponse{MessageType: messages.ProtocolMismatch})
	return SendTCPMessage(conn, messages.ProtocolMismatch, data, timeout)
}

// SendInvalidMessage sends an invalid message to a given connection
func SendInvalidMessage(conn net.Conn, timeout time.Duration, msg string) {
	data, _ := json.Marshal(messages.InvalidMessageResponse{MessageType: messages.InvalidMessage})
	logging.Info(msg)
	// Message is invalid.
	err := SendTCPMessage(conn, messages.InvalidMessage, data, timeout)
	if err != nil && !IsTimeoutError(err) {
		// Error in tcp communication, drop the connection.
		logging.Error1(err)
	}
}

// SendMessageWithType sends a given message with a given type to a given connection
func SendMessageWithType(conn net.Conn, msgType byte, v interface{}, timeout time.Duration) error {
	data, _ := json.Marshal(v)
	return SendTCPMessage(conn, msgType, data, timeout)
}
