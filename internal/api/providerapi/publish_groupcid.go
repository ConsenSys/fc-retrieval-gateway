package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
)

func handleProviderPublishGroupCIDRequest(conn net.Conn, request *messages.ProviderPublishGroupCIDRequest) error {
	g := gateway.GetSingleInstance()
	if g.Offers.Add(&cidoffer.CidGroupOffer{
		NodeID:               &request.ProviderID,
		Cids:                 request.PieceCIDs,
		Price:                request.Price,
		Expiry:               request.Expiry,
		QoS:                  request.QoS,
		Signature:            request.Signature,
		MerkleProof:          "TODO", //TODO: Who should be genearting the merkle proof?
		FundedPaymentChannel: true,   //TODO: Need to check if the gateway has a payment channel established with the provider.
	}) != nil {
		// Ignoved.
		logging.Error("Internal error in adding group cid offer.")
	}
	return nil
}
