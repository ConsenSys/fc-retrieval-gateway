package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

// StartProviderAPI starts the TCP API as a separate go routine.
func StartProviderAPI(settings settings.AppSettings) error {
	// Start server
	ln, err := net.Listen("tcp", ":"+settings.BindProviderAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			logging.Info("Incoming connection from provider at :%s\n", conn.RemoteAddr())
			go handleIncomingProviderConnection(conn)
		}
	}(ln)
	logging.Info("Listening on %s for connections from Providers\n", settings.BindProviderAPI)
	return nil
}

func handleIncomingProviderConnection(conn net.Conn) {
	// Close connection on exit.
	defer conn.Close()

	// Loop until error occurs and connection is dropped.
	for {
		message, err := fcrtcpcomms.ReadTCPMessage(conn, settings.DefaultTCPInactivityTimeout)

		if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error(err.Error())
			return
		}
		if err == nil {
			logging.Info("Message received: %+v", message)
			if message.MessageType == fcrmessages.ProviderPublishGroupCIDRequestType {
				handleProviderPublishGroupCIDRequest(message)
				continue
			} else if message.MessageType == fcrmessages.ProviderDHTPublishGroupCIDRequestType {
				err = handleProviderDHTPublishGroupCIDRequest(conn, message)
				if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection
					logging.Error(err.Error())
					return
				}
				continue
			}
			// Message is invalid.
			err = fcrtcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
			if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
				// Error in tcp communication, drop the connection.
				logging.Error(err.Error())
				return
			}
		}
	}
}
