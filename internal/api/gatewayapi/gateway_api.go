package gatewayapi

// Copyright (C) 2020 ConsenSys Software Inc

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

// StartGatewayAPI starts the TCP API as a separate go routine.
func StartGatewayAPI(settings settings.AppSettings) error {
	// Start server
	ln, err := net.Listen("tcp", ":"+settings.BindGatewayAPI)
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
			logging.Info("Incoming connection from gateway at :%s", conn.RemoteAddr())
			go handleIncomingGatewayConnection(conn)
		}
	}(ln)
	logging.Info("Listening on %s for connections from Gateways", settings.BindGatewayAPI)

	return nil
}

func handleIncomingGatewayConnection(conn net.Conn) {
	// Close connection on exit.
	defer conn.Close()

	// Loop until error occurs and connection is dropped.
	for {
		message, err := fcrtcpcomms.ReadTCPMessage(conn, settings.DefaultTCPInactivityTimeout)
		if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error1(err)
			return
		}
		if err == nil {
			if message.MessageType == fcrmessages.GatewayDHTDiscoverRequestType {
				err = handleGatewayDHTDiscoverRequest(conn, message)
				if err != nil && fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection
					logging.Error1(err)
					return
				}
				continue
			}
			// Message is invalid.
			err = fcrtcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
			if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
				// Error in tcp communication, drop the connection.
				logging.Error1(err)
				return
			}
		}
	}
}

// GetConnForRequestingGateway returns the connection for sending request to a gateway with given id.
// It will reuse any active connection.
func GetConnForRequestingGateway(gatewayID *nodeid.NodeID, g *gateway.Gateway) (*gateway.CommunicationChannel, error) {
	// Check if there is an active connection.
	g.ActiveGatewaysLock.RLock()
	gComm := g.ActiveGateways[gatewayID.ToString()]
	g.ActiveGatewaysLock.RUnlock()
	if gComm == nil {
		// No active connection, connect to peer.
		g.GatewayAddressMapLock.RLock()
		conn, err := net.Dial("tcp", g.GatewayAddressMap[gatewayID.ToString()])
		g.GatewayAddressMapLock.RUnlock()
		if err != nil {
			return nil, err
		}
		gComm = &gateway.CommunicationChannel{
			CommsLock: sync.RWMutex{},
			Conn:      conn}
		if gateway.RegisterGatewayCommunication(gatewayID, gComm) != nil {
			conn.Close()
			return nil, err
		}
	}
	return gComm, nil
}
