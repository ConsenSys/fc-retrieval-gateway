package adminapi

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"net"
	"sync"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

// StartAdminAPI starts the TCP API as a separate go routine.
func StartAdminAPI(settings settings.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", ":"+settings.BindAdminAPI)
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
			logging.Info("Incoming connection from admin client at :%s", conn.RemoteAddr())
			go handleIncomingAdminConnection(conn, g)
		}
	}(ln)
	logging.Info("Listening on %s for connections from admin clients", settings.BindAdminAPI)
	return nil
}

func handleIncomingAdminConnection(conn net.Conn, g *gateway.Gateway) {
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
		// Respond to requests for a client's reputation.
		if err == nil {
			if message.GetMessageType() == fcrmessages.GatewayAdminGetReputationRequestType {
				err = handleAdminGetReputationChallenge(conn, message)
				if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			} else if message.GetMessageType() == fcrmessages.GatewayAdminSetReputationRequestType {
				err = handleAdminSetReputationChallenge(conn, message)
				if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			} else if message.GetMessageType() == fcrmessages.GatewayAdminInitialiseKeyRequestType {
				var wg sync.WaitGroup
				wg.Add(1)
				err = handleAdminAcceptKeysChallenge(conn, message, &wg)
				if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection.
					logging.Error1(err)
					return
				}
				continue
			}
		}

		/*
			   TODO: Add additional message types:
			   ✔︎ Get a client's reputation.
			   - Set reputation of client arbitrarily.
			   - Set reputation of client based on various actions (e.g. using existing functionality).
			   - Set reputation of other gateway.
			   - Set reputation of provider.
			   - Get id of random client (for testing purposes).
			   - Remove Piece CID offers from the standard cache.
			   - Remove Piece CID offers from the DHT cache.
			   - Remove all Piece CID offers from a certain provider from the standard or DHT cache.
			   ✔︎ generate a key pair for the gateway.
				 - The API should have an optional parameter which
				is protocol version.
				- Store the private key in a runtime var (TODONEXT)
		*/

		// Message is invalid.
		fcrtcpcomms.SendInvalidMessage(conn, settings.DefaultTCPInactivityTimeout)
	}
}
