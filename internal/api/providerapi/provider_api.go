package providerapi

import (
	"net"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
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

// GetConnForRequestingProvider returns the connection for sending request to a provider with given id.
// It will reuse any active connection.
func GetConnForRequestingProvider(providerID *nodeid.NodeID, g *gateway.Gateway) (*gateway.CommunicationChannel, error) {
	// Check if there is an active connection.
	g.ActiveProvidersLock.RLock()
	pComm := g.ActiveProviders[providerID.ToString()]
	g.ActiveProvidersLock.RUnlock()
	if pComm == nil {
		// No active connection, connect to peer.
		g.ProviderAddressMapLock.RLock()
		conn, err := net.Dial("tcp", g.ProviderAddressMap[providerID.ToString()])
		g.ProviderAddressMapLock.RUnlock()
		if err != nil {
			return nil, err
		}
		pComm = &gateway.CommunicationChannel{
			CommsLock: sync.RWMutex{},
			Conn:      conn}
		if gateway.RegisterProviderCommunication(providerID, pComm) != nil {
			conn.Close()
			return nil, err
		}
	}
	return pComm, nil
}
