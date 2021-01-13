package providerapi

import (
	"encoding/json"
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

func handleProviderDHTPublishGroupCIDRequest(conn net.Conn, request *messages.ProviderDHTPublishGroupCIDRequest) error {
	g := gateway.GetSingleInstance()
	for _, offer := range request.CIDOffers {
		if g.Offers.Add(&cidoffer.CidGroupOffer{
			NodeID:               &request.ProviderID,
			Cids:                 []cid.ContentID{offer.PieceCID},
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
	sig, err := fcrcrypto.SignMessage(g.GatewayPrivateKey, g.GatewayPrivateKeyVersion, request)
	if err != nil {
		// Ignored.
		logging.Error("Error in signing message.")
	}

	// Respond to provider.
	response, _ := json.Marshal(messages.ProviderDHTPublishGroupCIDAck{
		MessageType:       messages.GatewayDHTDiscoverResponseType,
		ProtocolVersion:   g.ProtocolVersion,
		ProtocolSupported: g.ProtocolSupported,
		Nonce:             request.Nonce,
		Signature:         sig})
	return tcpcomms.SendTCPMessage(conn, messages.ProviderDHTPublishGroupCIDAckType, response, settings.DefaultTCPInactivityTimeout)
}
