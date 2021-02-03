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
	g.ProviderKeyMapLock.RLock()
	defer g.ProviderKeyMapLock.RUnlock()
	pubKey, ok := g.ProviderKeyMap[offer.NodeID.ToString()]
	if !ok {
		logging.Info("Provider public key not found.")
		return
	}

	ok, err = offer.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(&pubKey, sig, msg)
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
