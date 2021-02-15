package providerapi

import (
	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

func handleProviderPublishGroupCIDRequest(request *fcrmessages.FCRMessage) {
	// Get the core structure
	g := gateway.GetSingleInstance()

	// TODO: Why we need a nonce here?
	_, offer, err := fcrmessages.DecodeProviderPublishGroupCIDRequest(request)
	if err != nil {
		logging.Info("Provider publish group cid request fail to decode.")
		return
	}

	// Need to verify the offer
	// Get the public key
	g.RegisteredProvidersMapLock.RLock()
	defer g.RegisteredProvidersMapLock.RUnlock()
	provider, ok := g.RegisteredProvidersMap[offer.NodeID.ToString()]
	if !ok {
		logging.Info("Provider public key not found.")
		return
	}
	pubKey, err := provider.GetSigningKey()
	if err != nil {
		logging.Info("Fail to get signing key from provider registration info")
		return
	}

	ok, err = offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})

	// Error in verifying message
	if err != nil {
		logging.Error("Internal error in verifying group cid offer.")
		// Ignored.
		return
	}

	// Offer does not pass verification
	if !ok {
		logging.Info("Offer does not pass verification.")
		// Ignored.
		return
	}

	if g.Offers.Add(offer) != nil {
		// Ignoved.
		logging.Error("Internal error in adding group cid offer.")
	}
	return
}
