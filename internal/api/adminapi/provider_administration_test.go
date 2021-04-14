package adminapi

import (
	"testing"

	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
)

func TestRegisterProvider_ProviderAddedWhenDoesNotExist(t *testing.T) {
	// arrange
	expectedNodeID := "correct_node_id"
	newProvider := register.ProviderRegister{
		NodeID: expectedNodeID,
	}
	gatewayCoreStructure := gateway.Gateway{
		RegisteredProvidersMap: map[string]register.RegisteredNode{},
	}

	// act
	registerProvider(newProvider, &gatewayCoreStructure)

	// assert
	_, exists := gatewayCoreStructure.RegisteredProvidersMap[expectedNodeID]
	if !exists {
		t.Error("A new providers can't be registered")
	}
}
