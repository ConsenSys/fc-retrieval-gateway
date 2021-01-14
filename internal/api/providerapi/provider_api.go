package providerapi

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// StartProviderAPI starts the TCP API as a separate go routine.
func StartProviderAPI(settings settings.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", settings.BindProviderAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logging.Error1(err)
				continue
			}
			logging.Info("Incoming connection from provider at :%s\n", conn.RemoteAddr())
			go handleIncomingProviderConnection(conn, g)
		}
	}(ln)
	logging.Info("Listening on %s for connections from Providers\n", settings.BindProviderAPI)
	return nil
}

func handleIncomingProviderConnection(conn net.Conn, g *gateway.Gateway) {
	// Close connection on exit.
	defer conn.Close()

	// Loop until error occurs and connection is dropped.
	for {
		msgType, data, err := tcpcomms.ReadTCPMessage(conn, settings.DefaultTCPInactivityTimeout)
		if err != nil && !tcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error1(err)
			return
		}
		if msgType == messages.ProviderPublishGroupCIDRequestType {
			request := messages.ProviderPublishGroupCIDRequest{}
			if json.Unmarshal(data, &request) == nil {
				// Message is valid.
				handleProviderPublishGroupCIDRequest(conn, &request)
				continue
			}
		} else if msgType == messages.ProviderDHTPublishGroupCIDRequestType {
			request := messages.ProviderDHTPublishGroupCIDRequest{}
			if json.Unmarshal(data, &request) == nil {
				// Message is valid.
				err = handleProviderDHTPublishGroupCIDRequest(conn, &request)
				if err != nil && !tcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			}
		}
		// Message is invalid.
		err = tcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
		if err != nil && !tcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error1(err)
			return
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
