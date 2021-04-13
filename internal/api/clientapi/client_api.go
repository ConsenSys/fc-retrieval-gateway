package clientapi

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
	"io/ioutil"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

// StartClientRestAPI starts the REST API as a separate go routine.
// Any start-up errors are returned.
func StartClientRestAPI(settings settings.AppSettings) error {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/v1", msgRouter),
	)
	if err != nil {
		return err
	}
	api.SetApp(router)
	err = http.ListenAndServe(":"+settings.BindRestAPI, api.MakeHandler())
	if err != nil {
		return err
	}
	return nil
}

// msgRouter routes message
func msgRouter(w rest.ResponseWriter, r *rest.Request) {
	// Get core structure
	c := core.GetSingleInstance()

	logging.Trace("Received request via /v1 API")
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		logging.Error("Error reading request: %s.", err.Error())
		rest.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	if len(content) == 0 {
		logging.Error("Error empty request")
		rest.Error(w, "Error empty request", http.StatusBadRequest)
		return
	}
	request, err := fcrmessages.FCRMsgFromBytes(content)
	if err != nil {
		logging.Error("Failed to decode payload: %s.", err.Error())
		rest.Error(w, "Failed to decode payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Only process the rest of the message if the protocol version is understood.
	if request.GetProtocolVersion() != c.ProtocolVersion {
		// Check to see if the client supports the gateway's preferred version
		for _, clientProvVer := range request.GetProtocolSupported() {
			if clientProvVer == c.ProtocolVersion {
				// Request the client switch to this protocol version
				// TODO what can we get from request object?
				logging.Info("Requesting client (TODO) switch protocol versions from %d to %d", request.GetProtocolVersion(), c.ProtocolVersion)
				response, _ := fcrmessages.EncodeProtocolChangeRequest(c.ProtocolVersion)
				w.WriteJson(response)
				return
			}
		}

		// Go through the protocol versions supported by the client and the
		// gateway to search for any common version, prioritising
		// the gateway preference over the client preference.
		for _, clientProvVer := range request.GetProtocolSupported() {
			for _, gatewayProtVer := range c.ProtocolSupported {
				if clientProvVer == gatewayProtVer {
					// When we support more than one version of the protocol, this code will change the gateway
					// to using the other (common version)
					logging.Error("Not implemented yet")
					panic("Multiple protocol versions not implemented yet")
				}
			}
		}
		// No common protocol versions supported.
		// TODO what can we get from request object?
		logging.Warn("Client Request: Unsupported protocol version(s): %d", request.GetProtocolVersion())
		response, _ := fcrmessages.EncodeProtocolChangeResponse(false)
		w.WriteJson(response)
		return
	}

	switch request.GetMessageType() {
	case fcrmessages.ClientEstablishmentRequestType:
		handleClientEstablishmentRequest(w, request)
	case fcrmessages.ClientStandardDiscoverRequestType:
		handleClientStandardCIDDiscoverRequest(w, request)
	case fcrmessages.ClientDHTDiscoverRequestType:
		handleClientDHTCIDDiscoverRequest(w, request)
	default:
		logging.Warn("Client Request: Unknown message type: %d", request.GetMessageType())
		rest.Error(w, "Unknown message type", http.StatusBadRequest)
	}
}
