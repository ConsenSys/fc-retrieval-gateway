package dht

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
	"fmt"

	"github.com/ConsenSys/fc-retrieval-common/pkg/dhtring"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/core"
)

func UpdateNearestGatewaysDHT() {
	// gateway instance.
	c := core.GetSingleInstance()

	logging.Info("---------------------- UpdateNearestGatewaysDHT Start----------------------")
	defer logging.Info("---------------------- UpdateNearestGatewaysDHT End ----------------------")

	allGateways := c.RegisterMgr.GetAllGateways()
	if len(allGateways) == 0 || c.GatewayID == nil {
		return
	}

	gtwIDs := make([]*nodeid.NodeID, 0)
	contacted := make([]fcrmessages.FCRMessage, 0)

	for _, gtw := range allGateways {
		nodeID := gtw.NodeID
		if c.GatewayID.ToString() == nodeID {
			continue
		}

		gtwID, err := nodeid.NewNodeIDFromHexString(nodeID)
		if err != nil {
			logging.Error("error getting nodeID %s", nodeID)
			continue
		}

		res, err := c.P2PServer.RequestGatewayFromGateway(gtwID, fcrmessages.GatewayPingRequestType, gtwID)
		if err != nil {
			logging.Error("gatewayID not available error: %s, %s", err, nodeID)
			continue
		}
		logging.Debug("gatewayID available ! %s", nodeID)

		contacted = append(contacted, *res)

		fmt.Printf("%v\n", contacted)

		gtwIDs = append(gtwIDs, gtwID)
	}

	closestGtwIDs, err := dhtring.SortClosestNodesIDs(c.GatewayID.ToBytes(), gtwIDs)
	if err != nil {
		logging.Error("Cant found closest allGateways.")
	}

	closestGatewaysIDs := make([][]byte, 0)
	for _, gtw := range closestGtwIDs {
		closestGatewaysIDs = append(closestGatewaysIDs, gtw.ToBytes())
	}

	logging.Debug("c.RegisterMgr.SetClosestGatewaysIDs %v", closestGatewaysIDs)

	c.RegisterMgr.SetClosestGatewaysIDs(closestGatewaysIDs)
}
