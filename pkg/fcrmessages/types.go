package fcrmessages

import "github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// Message types
// The enum should remains the same for client,provider and gateway.
const (
	ClientEstablishmentRequestType              = 0
	ClientEstablishmentResponseType             = 1
	ClientStandardDiscoverRequestType           = 2
	ClientStandardDiscoverResponseType          = 3
	ClientDHTDiscoverRequestType                = 4
	ClientDHTDiscoverResponseType               = 5
	ClientCIDGroupPublishDHTAckRequestType      = 6
	ClientCIDGroupPublishDHTAckResponseType     = 7
	ProviderPublishGroupCIDRequestType          = 8
	ProviderDHTPublishGroupCIDRequestType       = 9
	ProviderDHTPublishGroupCIDAckType           = 10
	GatewaySingleCIDOfferPublishRequestType     = 11
	GatewaySingleCIDOfferPublishResponseType    = 12
	GatewaySingleCIDOfferPublishResponseAckType = 13
	GatewayDHTDiscoverRequestType               = 14
	GatewayDHTDiscoverResponseType              = 15
	ProtocolChangeResponseType                  = 100
	ProtocolMismatchResposneType                = 101
	InvalidMessageResponseType                  = 102
	InsufficientFundsResponseType               = 103
	AdminGetReputationChallengeType             = 200
	AdminGetReputationResponseType              = 201
	AdminSetReputationChallengeType             = 202
	AdminSetReputationResponseType              = 203
)

// CIDGroupInformation represents a cid group information
type CIDGroupInformation struct {
	ProviderID           nodeid.NodeID `json:"provider_id"`
	Price                uint64        `json:"price_per_byte"`
	Expiry               int64         `json:"expiry_date"`
	QoS                  uint64        `json:"qos"`
	Signature            string        `json:"signature"`
	MerkleProof          string        `json:"merkle_proof"`
	FundedPaymentChannel bool          `json:"funded_payment_channel"` // TODO: Is this boolean?
}
