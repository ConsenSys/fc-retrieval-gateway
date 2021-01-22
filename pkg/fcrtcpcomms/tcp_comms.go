package fcrtcpcomms

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
)

// IsTimeoutError checks if the given error is a timeout error
func IsTimeoutError(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr.Timeout()
}

// ReadTCPMessage read the tcp message from a given connection
func ReadTCPMessage(conn net.Conn, timeout time.Duration) (*fcrmessages.FCRMessage, error) {
	// Initialise a reader
	reader := bufio.NewReader(conn)
	// Read the length
	length := make([]byte, 4)
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	_, err := io.ReadFull(reader, length)
	if err != nil {
		return nil, err
	}
	// Read the data
	data := make([]byte, int(binary.BigEndian.Uint32(length)))
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	return fcrmessages.FCRMsgFromBytes(data)
}

// SendTCPMessage sends a tcp message to a given connection
func SendTCPMessage(conn net.Conn, fcrMsg *fcrmessages.FCRMessage, timeout time.Duration) error {
	// Get data
	data, err := fcrmessages.FCRMsgToBytes(fcrMsg)
	if err != nil {
		return err
	}
	// Initialise a writer
	writer := bufio.NewWriter(conn)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(data)))
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	_, err = writer.Write(append(length, data...))
	if err != nil {
		return err
	}
	// Set timeout
	conn.SetDeadline(time.Now().Add(timeout))
	return writer.Flush()
}

// SendProtocolMismatch sends a protocol mistmatch message to a given connection
func SendProtocolMismatch(conn net.Conn, timeout time.Duration) error {
	fcrMsg, _ := fcrmessages.EncodeProtocolMismatchResponse()
	return SendTCPMessage(conn, fcrMsg, timeout)
}

// SendInvalidMessage sends an invalid message to a given connection
func SendInvalidMessage(conn net.Conn, timeout time.Duration) error {
	fcrMsg, _ := fcrmessages.EncodeInvalidMessageResponse()
	return SendTCPMessage(conn, fcrMsg, timeout)
}
