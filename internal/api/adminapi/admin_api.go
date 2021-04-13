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
	"io/ioutil"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ant0ine/go-json-rest/rest"
)

// StartAdminAPI starts the TCP API as a separate go routine.
func StartAdminAPI(settings settings.AppSettings) error {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/v1", msgRouter),
	)
	if err != nil {
		return err
	}
	api.SetApp(router)
	err = http.ListenAndServe(":"+settings.BindAdminAPI, api.MakeHandler())
	if err != nil {
		return err
	}
	return nil
}

// msgRouter routes message
func msgRouter(w rest.ResponseWriter, r *rest.Request) {
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

	switch request.GetMessageType() {
	case fcrmessages.GatewayAdminInitialiseKeyRequestType:
		handleGatewayAdminInitialiseKeyRequest(w, request)
	case fcrmessages.GatewayAdminGetReputationRequestType:
		handleGatewayAdminGetReputationRequest(w, request)
	case fcrmessages.GatewayAdminSetReputationRequestType:
		handleGatewayAdminSetReputationRequest(w, request)
	default:
		logging.Warn("Client Request: Unknown message type: %d", request.GetMessageType())
		rest.Error(w, "Unknown message type", http.StatusBadRequest)
	}
}
