package providerapi

import (
	// "errors"
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"

	// "github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

func handleProviderPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	logging.Info("handleProviderPublishGroupCIDRequest: %+v", request)

	// Get the core structure
	g := gateway.GetSingleInstance()

	// TODO: Why we need a nonce here?
	logging.Info("Decode provider publish group CID request: %+v", request)
	_, offer, err := fcrmessages.DecodeProviderPublishGroupCIDRequest(request)
	if err != nil {
		logging.Info("Provider publish group cid request fail to decode.")
		return err
	}

	// Need to verify the offer
	// Get the public key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	// provider, ok := g.RegisteredProvidersMap[offer.NodeID.ToString()]
	// if !ok {
	// 	logging.Info("Provider public key not found.")
	// 	return errors.New("Provider public key not found")
	// }
	// pubKey, err := provider.GetSigningKey()
	// if err != nil {
	// 	logging.Info("Fail to get signing key from provider registration info")
	// 	return err
	// }

	// ok, err = offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
	// 	return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	// })

	// Error in verifying message
	if err != nil {
		logging.Error("Internal error in verifying group cid offer.")
		// Ignored.
		return err
	}

	// Offer does not pass verification
	// if !ok {
	// 	logging.Info("Offer does not pass verification.")
	// 	// Ignored.
	// 	return errors.New("Offer does not pass verification")
	// }

	// Store the offer
	if g.Offers.Add(offer) != nil {
		logging.Error("Internal error in adding group cid offer.")
		return err
	}

	logging.Info("Encode provider publish group CID response: %+v", offer)
	response, err := fcrmessages.EncodeProviderPublishGroupCIDResponse(
		*g.GatewayID,
		offer.GetMessageDigest(),
	)
	if err != nil {
		logging.Error("Internal error in encoding publish group cid response.")
		return err
	}

	logging.Info("Send response to provider: %+v", response)
	return fcrtcpcomms.SendTCPMessage(conn, response, settings.DefaultTCPInactivityTimeout)
}
