package providerapi

import (
	"errors"
	"net"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

func handleProviderDHTPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	providerID, nonce, offers, err := fcrmessages.DecodeProviderPublishDHTOfferRequest(request)
	if err != nil {
		logging.Info("Provider publish dht offer request fail to decode.")
		return nil
	}

	// Get the provider's signing key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	_, ok := g.RegisteredProvidersMap[providerID.ToString()]
	if !ok {
		return errors.New("Provider public key not found")
	}
	pubKey, err := g.RegisteredProvidersMap[providerID.ToString()].GetSigningKey()
	if err != nil {
		return err
	}
	// First verify the message
	if request.Verify(pubKey) != nil {
		return errors.New("Fail to verify the request")
	}

	// Second Need to verify the offer one by one
	for _, offer := range offers {
		if offer.Verify(pubKey) != nil {
			return errors.New("Fail to verify the offer")
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

	response, err := fcrmessages.EncodeProviderPublishDHTOfferResponse(nonce, sig)
	if err != nil {
		logging.Error("Internal error in encoding message.")
		return nil
	}
	// Sign the response
	if response.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		return errors.New("Error in signing message")
	}
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
