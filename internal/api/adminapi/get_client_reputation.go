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
	"net/http"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/reputation"
	"github.com/ant0ine/go-json-rest/rest"
)

// handleAdminGetReputationChallenge handles admin get reputation challenge
func handleAdminGetReputationChallenge(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	if c.GatewayPrivateKey == nil {
		s := "This gateway hasn't been initialised by the admin"
		logging.Error(s)
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	clientID, err := fcrmessages.DecodeGatewayAdminGetReputationRequest(request)
	if err != nil {
		s := "Admin Get Reputation: Failed to decode payload."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	// Get reputation db
	rep := reputation.GetSingleInstance()
	reputation, exists := rep.GetClientReputation(clientID)

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminGetReputationResponse(clientID, reputation, exists)
	if err != nil {
		s := "Internal error: Fail to encode response."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	// Sign message
	if response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion) != nil {
		s := "Internal error."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	w.WriteJson(response)
}
