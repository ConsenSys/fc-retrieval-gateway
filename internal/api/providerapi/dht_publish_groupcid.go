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

	nonce, providerID, offers, err := fcrmessages.DecodeProviderDHTPublishGroupCIDRequest(request)
	if err != nil {
		logging.Info("Provider publish cid offer dht request fail to decode.")
		return nil
	}

	// Get the public key
	g.ProviderKeyMapLock.RLock()
	defer g.ProviderKeyMapLock.RUnlock()
	pubKey, ok := g.ProviderKeyMap[providerID.ToString()]
	if !ok {
		logging.Info("Provider public key not found.")
		return nil
	}

	for _, offer := range offers {
		// Need to verify the offer one by one
		ok, err = offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
			return fcrcrypto.VerifyMessage(&pubKey, sig, msg)
		})

		if err != nil {
			logging.Error("Internal error in verifying the cid offer.")
			continue
		}

		if !ok {
			logging.Info("Offer does not pass verification.")
			continue
		}

		if g.Offers.Add(&offer) != nil {
			// Ignored.
			logging.Error("Internal error in adding single cid offer.")
		}
	}

	// Sign the request
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, request)
	if err != nil {
		// Ignored.
		logging.Error("Internal error in signing message.")
		return nil
	}

	response, err := fcrmessages.EncodeProviderDHTPublishGroupCIDAck(nonce, sig)
	if err != nil {
		logging.Error("Internal error in encoding message.")
		return nil
	}

	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
