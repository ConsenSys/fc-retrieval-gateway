/*
Package reputation - contains operations related to blockchain network participants' favorability marking.
Participant with higher reputation is generally more preferable.
*/
package reputation

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

import (
	"sync"

	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// Single instance of the reputation system
var instance = newInstance()

// Reputation manages the reputation of all other nodes in the system from this
// node's perspective.
type Reputation struct {
	clients          map[string]int64
	clientsMapLock   sync.RWMutex
	gateways         map[string]int64
	gatewaysMapLock  sync.RWMutex
	providers        map[string]int64
	providersMapLock sync.RWMutex
}

// Create a new instance
func newInstance() *Reputation {
	var r = Reputation{}
	r.clients = make(map[string]int64)
	r.clientsMapLock = sync.RWMutex{}
	r.gateways = make(map[string]int64)
	r.gatewaysMapLock = sync.RWMutex{}
	r.providers = make(map[string]int64)
	r.providersMapLock = sync.RWMutex{}
	return &r
}

// GetSingleInstance is a factory method to get the single instance of the reputation system
func GetSingleInstance() *Reputation {
	return instance
}

// establishClientReputation initialises the reputation of a Retrieval Client
func (r *Reputation) establishClientReputation(clientNodeID *nodeid.NodeID) {
	r.setClientReputation(clientNodeID, clientInitialReputation)
}

// SetClientReputation sets the reputation of a Retrieval Client
func (r *Reputation) SetClientReputation(clientNodeID *nodeid.NodeID, newReputation int64) {
	r.setClientReputation(clientNodeID, newReputation)
}

// ClientExists determines if a client has reputation
func (r *Reputation) ClientExists(clientNodeID *nodeid.NodeID) (exists bool) {
	_, exists = r.getClientReputation(clientNodeID)
	return
}

// GetClientReputation returns the client reputation.
func (r *Reputation) GetClientReputation(clientNodeID *nodeid.NodeID) (val int64, exists bool) {
	val, exists = r.getClientReputation(clientNodeID)
	return
}

// ClientEstablishmentChallenge updates a Retrieval Client's reputation based on an
// Establishment Challenge being received. The reputation is created for the Retrival
// Client if the client doesn't have a reputation yet
func (r *Reputation) ClientEstablishmentChallenge(clientNodeID *nodeid.NodeID) int64 {
	_, exists := r.getClientReputation(clientNodeID)
	if !exists {
		r.establishClientReputation(clientNodeID)
	}
	return r.changeClientReputation(clientNodeID, clientEstablishmentChallenge)
}

// OnChainDeposit updates a Retrieval Client's reputation based on an
// a deposit of Establishment Challenge being received.
func (r *Reputation) OnChainDeposit(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientOnChainDeposit)
}

// ClientStdDiscOneCidOffer updates reputation given a response with one or more CID Offers.
// Initial payment and final payment made.
func (r *Reputation) ClientStdDiscOneCidOffer(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientStdDiscOneCidOffer)
}

// ClientStdDiscNoCidOffers updates reputation given a response with no CID Offers. Initial
// payment payment made.
func (r *Reputation) ClientStdDiscNoCidOffers(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientStdDiscNoCidOffers)
}

// ClientStdDiscLateCidOffers updates reputation given a response with one or more CID Offers.
// Response message sent after one second prior to TTL expiry. Initial payment payment made.
func (r *Reputation) ClientStdDiscLateCidOffers(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientStdDiscLateCidOffers)
}

// ClientStdDiscNonPayment updates reputation given a response with one or more CID Offers.
// Response message sent prior to one second prior to TTL expiry. Initial payment payment
// made but final payment not paid.
func (r *Reputation) ClientStdDiscNonPayment(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientStdDiscNonPayment)
}

// ClientDhtDiscOneCidOffer updates reputation given a response with one or more CID Offers
// from one or more Gateways. Initial payment and final payment made
func (r *Reputation) ClientDhtDiscOneCidOffer(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientDhtDiscOneCidOffer)
}

// ClientDhtDiscNoCidOffers updates reputation given a response with no CID Offers.
// Initial payment made.
func (r *Reputation) ClientDhtDiscNoCidOffers(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientDhtDiscNoCidOffers)
}

// ClientDhtDiscLateCidOffers updates reputation given a response with one or more CID
// Offers. Response message sent after one second prior to TTL expiry. Initial payment
// payment made.
func (r *Reputation) ClientDhtDiscLateCidOffers(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientDhtDiscLateCidOffers)
}

// ClientDhtDiscNonPayment updates reputation given a response with one or more CID
// Offers. Response message sent prior to one second prior to TTL expiry. Initial payment
// payment made but final payments not paid.
func (r *Reputation) ClientDhtDiscNonPayment(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientDhtDiscNonPayment)
}

// ClientMicroPayment updates reputation given a micro-payment paid for content via
// Gateway. Note that there will be many micro-payments during content retrieval.
func (r *Reputation) ClientMicroPayment(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientMicroPayment)
}

// ClientInvalidMessage updates reputation given an invalid message received.
func (r *Reputation) ClientInvalidMessage(clientNodeID *nodeid.NodeID) {
	r.changeClientReputation(clientNodeID, clientInvalidMessage)
}
