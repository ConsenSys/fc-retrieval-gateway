package messages

// Message types
// The enum should remains the same for client,provider and gateway.
const (
	ClientEstablishmentRequestType = 0
	ClientEstablishmentResponseType = 1
	ClientStandardDiscoverRequestType = 2
	ClientStandardDiscoverResponseType = 3
	ClientDHTDiscoverRequestType = 4
	ClientDHTDiscoverResponseType = 5
	ProtocolChange = 100
	ProtocolMismatch = 101
)