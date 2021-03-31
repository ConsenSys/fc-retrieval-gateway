package providerapi

import (
	"errors"
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

func handleProviderPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// TODO Add nonce, it looks like nonce is not needed
	providerID, _, offer, err := fcrmessages.DecodeProviderPublishGroupOfferRequest(request)
	if err != nil {
		logging.Info("Provider publish group cid request fail to decode.")
		return err
	}
	// Get the provider's signing key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	_, ok := g.RegisteredProvidersMap[providerID.ToString()]
	if !ok {
		return errors.New("Provider not found")
	}
	pubKey, err := g.RegisteredProvidersMap[providerID.ToString()].GetSigningKey()
	if err != nil {
		return err
	}
	// First verify the message
	if request.Verify(pubKey) != nil {
		return errors.New("Fail to verify the request")
	}

	// Second Need to verify the offer
	if offer.Verify(pubKey) != nil {
		return errors.New("Fail to verify the offer")
	}

	// Store the offer
	if g.Offers.Add(offer) != nil {
		return errors.New("Fail to store the offer")
	}
	logging.Info("Stored offers: %+v", g.Offers)

	response, err := fcrmessages.EncodeProviderPublishGroupOfferResponse(
		*g.GatewayID,
		offer.GetMessageDigest(),
	)
	if err != nil {
		logging.Error("Internal error in encoding publish group cid response.")
		return err
	}

	// Sign the response
	if response.Sign(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion) != nil {
		return errors.New("Error in signing message")
	}
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
