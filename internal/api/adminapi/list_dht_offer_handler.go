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

  "github.com/ant0ine/go-json-rest/rest"

  "github.com/ConsenSys/fc-retrieval-common/pkg/cid"
  "github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
  "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
  "github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
  "github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

// HandleGatewayAdminListDHTOffersRequest handles admin list dht offer request
func HandleGatewayAdminListDHTOffersRequest(w rest.ResponseWriter, request *fcrmessages.FCRMessage) {
	// Get core structure
	c := core.GetSingleInstance()

	if c.GatewayPrivateKey == nil {
		s := "This gateway hasn't been initialised by the admin"
		logging.Error(s)
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	refresh, err := fcrmessages.DecodeGatewayAdminListDHTOfferRequest(request)
	if err != nil {
		s := "Fail to decode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusBadRequest)
		return
	}

	refreshed := false
	if refresh {
		// Send list cid offers
		go func() {
			cidMin, cidMax, err := c.RegisterMgr.GetGatewayCIDRange(c.GatewayID)
			if err != nil {
				logging.Error("Error getting cid range for gateway initial list dht offer: %v", err.Error())
				return
			}
			pvds := c.RegisterMgr.GetAllProviders()
			for _, pvd := range pvds {
				id, err := nodeid.NewNodeIDFromHexString(pvd.GetNodeID())
				if err != nil {
					logging.Error("Error in generating node id")
					continue
				}
				go requestProvider(c, id, fcrmessages.GatewayListDHTOfferRequestType, cidMin, cidMax, id)
			}
		}()
		refreshed = true
	}

	// Construct message
	response, err := fcrmessages.EncodeGatewayAdminListDHTOfferResponse(refreshed)
	if err != nil {
		s := "Internal error: Fail to encode message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
	// Sign message
	err = response.Sign(c.GatewayPrivateKey, c.GatewayPrivateKeyVersion)
	if err != nil {
		s := "Internal error: Fail to sign message."
		logging.Error(s + err.Error())
		rest.Error(w, s, http.StatusInternalServerError)
		return
	}
  if err := w.WriteJson(response); err != nil {
    logging.Error("can't write JSON during HandleGatewayAdminListDHTOffersRequest %s", err.Error())
  }
}

func requestProvider(gatewayInstance *core.Core, providerID *nodeid.NodeID, msgType int32, cidMin *cid.ContentID, cidMax *cid.ContentID, nodeID *nodeid.NodeID) {
  providerResponse, err := gatewayInstance.P2PServer.RequestProvider(providerID, msgType, cidMin, cidMax, nodeID)
  if err != nil {
    logging.Error("error requesting provider, request message type: %d; provider id: %s, error: %s", msgType, providerID.ToString(), err.Error())
    return
  }
  if providerResponse == nil {
    logging.Warn("gateway can't request a provider id: %s, request message type: %d; P2P method RequestProvider returned nil", providerID.ToString(), msgType)
    return
  }
  logging.Debug("Provider response: %s for request message type: %d", providerResponse.DumpMessage(), msgType)
}
