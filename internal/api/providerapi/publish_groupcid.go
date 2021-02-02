package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

func handleProviderPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	_, offer, err := fcrmessages.DecodeProviderPublishGroupCIDRequest(request)
	if err != nil {
		return err
	}

	// TODO: Need to verify each offer by signature

	if g.Offers.Add(offer) != nil {
		// Ignoved.
		logging.Error("Internal error in adding group cid offer.")
	}
	return nil
}
