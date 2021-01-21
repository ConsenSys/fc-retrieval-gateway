package messages

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
	ProviderPublishGroupCIDRequestType          = 6
	ProviderDHTPublishGroupCIDRequestType       = 7
	ProviderDHTPublishGroupCIDAckType           = 8
	GatewaySingleCIDOfferPublishRequestType     = 9
	GatewaySingleCIDOfferPublishResponseType    = 10
	GatewaySingleCIDOfferPublishResponseAckType = 11
	GatewayDHTDiscoverRequestType               = 12
	GatewayDHTDiscoverResponseType              = 13
	ProtocolChange                              = 100
	ProtocolMismatch                            = 101
	InvalidMessage                              = 102
	InsufficientFunds                           = 103
	AdminGetReputationChallengeType             = 200
	AdminGetReputationResponseType              = 201
	AdminSetReputationChallengeType             = 202
	AdminSetReputationResponseType              = 203
)
