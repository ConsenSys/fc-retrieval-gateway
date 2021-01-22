package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
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
		return err
	}

	for _, offer := range offers {
		if g.Offers.Add(&cidoffer.CidGroupOffer{
			NodeID:               providerID,
			Cids:                 offer.Cids,
			Price:                offer.Price,
			Expiry:               offer.Expiry,
			QoS:                  offer.QoS,
			Signature:            offer.Signature,
			MerkleProof:          "TODO", //TODO: Who should be genearting the merkle proof?
			FundedPaymentChannel: true,   //TODO: Need to check if the gateway has a payment channel established with the provider.
		}) != nil {
			// Ignored.
			logging.Error("Internal error in adding single cid offer.")
		}
	}
	// Sign the message
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, request.MessageBody)
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
