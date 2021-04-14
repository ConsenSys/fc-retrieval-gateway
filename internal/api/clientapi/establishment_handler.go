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
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ant0ine/go-json-rest/rest"
)

// HandleClientEstablishmentRequest is used to handle initial establishment http request from client
func HandleClientEstablishmentRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	clientID, challenge, ttl, err := fcrmessages.DecodeClientEstablishmentRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	logging.Trace("Client Establishment from %s with challenge %s and ttl %d", clientID.ToString(), challenge, ttl)

	now := util.GetTimeImpl().Now().Unix()
	if now > ttl {
		// TODO how to just drop the connection?
		return
	}

	// Construct message
	response, err := fcrmessages.EncodeClientEstablishmentResponse(c.GatewayID, challenge)
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// Sign message
	if response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	w.WriteJson(response)
}
