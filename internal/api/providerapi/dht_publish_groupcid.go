package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

func handleProviderDHTPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	nonce, _, offers, err := fcrmessages.DecodeProviderDHTPublishGroupCIDRequest(request)
	if err != nil {
		return err
	}

	for _, offer := range offers {
		// TODO: Need to verify each offer by signature
		if g.Offers.Add(&offer) != nil {
			// Ignored.
			logging.Error("Internal error in adding single cid offer.")
		}
	}
	// Sign the message
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, request)
	if err != nil {
		// Ignored.
		logging.Error("Error in signing message.")
	}

	response, err := fcrmessages.EncodeProviderDHTPublishGroupCIDAck(nonce, sig)
	if err != nil {
		return err
	}

	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
