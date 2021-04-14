package adminapi

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
)

func TestRegisterGateway_GatewayAddedWhenDoesNotExist(t *testing.T) {
	// arrange
	expectedNodeID := "correct_node_id"
	newGateway := register.GatewayRegister{
		NodeID: expectedNodeID,
	}
	gatewayCoreStructure := gateway.Gateway{
		RegisteredGatewaysMap: map[string]register.RegisteredNode{},
	}

	// act
	registerGateway(newGateway, &gatewayCoreStructure)

	// assert
	_, exists := gatewayCoreStructure.RegisteredGatewaysMap[expectedNodeID]
	if !exists {
		t.Error("A new gateway can't be registered")
	}
}
