package dht

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

	logging.Info("-------------------------------------------------------------")

	allGateways := c.RegisterMgr.GetAllGateways()
	if len(allGateways) == 0 || c.Settings.GatewayID == "" {
		return
	}

	gtwIDs := make([]*nodeid.NodeID, 0)
	contacted := make([]fcrmessages.FCRMessage, 0)

	for _, gtw := range allGateways {
		gtwID, err := nodeid.NewNodeIDFromHexString(gtw.GetNodeID())
		if err != nil {
			logging.Error("error getting nodeID %s", gtw.GetNodeID())
			continue
		}

		res, err := c.P2PServer.RequestGatewayFromGateway(gtwID, fcrmessages.GatewayPingRequestType)
		if err != nil {
			logging.Error("gatewayID not available %s", gtw.GetNodeID())
			continue
		}

		logging.Info("---------------", gtwIDs, contacted, res)

		contacted = append(contacted, *res)
		// res.
		fmt.Printf("%v\n", contacted)

		gtwIDs = append(gtwIDs, gtwID)
	}

	closestGtwIDs, err := dhtring.SortClosestNodesIDs(c.Settings.GatewayID, gtwIDs)
	if err != nil {
		logging.Error("Cant found closest allGateways.")
	}

	closestGatewaysIDs := make([]string, 0)
	for _, gtw := range closestGtwIDs {
		closestGatewaysIDs = append(closestGatewaysIDs, gtw.ToString())
	}

	c.RegisterMgr.SetClosestGatewaysIDs(closestGatewaysIDs)

}
