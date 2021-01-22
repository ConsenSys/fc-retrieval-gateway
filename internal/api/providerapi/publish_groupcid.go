package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

func handleProviderPublishGroupCIDRequest(conn net.Conn, request *fcrmessages.FCRMessage) error {
	// Get the core structure
	g := gateway.GetSingleInstance()

	_, providerID, price, expiry, qos, cids, err := fcrmessages.DecodeProviderPublishGroupCIDRequest(request)
	if err != nil {
		return err
	}

	if g.Offers.Add(&cidoffer.CidGroupOffer{
		NodeID:               providerID,
		Cids:                 cids,
		Price:                price,
		Expiry:               expiry,
		QoS:                  qos,
		Signature:            request.Signature,
		MerkleProof:          "TODO", //TODO: Who should be genearting the merkle proof?
		FundedPaymentChannel: true,   //TODO: Need to check if the gateway has a payment channel established with the provider.
	}) != nil {
		// Ignoved.
		logging.Error("Internal error in adding group cid offer.")
	}
	return nil
}
