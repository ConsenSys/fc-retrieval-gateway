package core

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

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmerkletree"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcroffermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrp2pserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrrestserver"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/reputation"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

const (
	protocolVersion   = 1 // Main protocol version
	protocolSupported = 1 // Alternative protocol version
)

// Core holds the main data structure for the whole gateway.
type Core struct {
	// Protocol versions of this gateway
	ProtocolVersion   int32
	ProtocolSupported []int32

	// Settings
	Settings *settings.AppSettings

	// GatewayID of this gateway
	GatewayID *nodeid.NodeID

	// Gateway Private Key and algorithm of this gateway
	GatewayPrivateKey *fcrcrypto.KeyPair

	// GatewayPrivateKeyVersion is the key version number of the private key.
	GatewayPrivateKeyVersion *fcrcrypto.KeyVersion

	// RegisterMgr manages all register related activities
	RegisterMgr *fcrregistermgr.FCRRegisterMgr

	// P2PServer handles all communication to/from gateways/providers
	P2PServer *fcrp2pserver.FCRP2PServer

	// RESTServer handles all communication to/from client/admin
	RESTServer *fcrrestserver.FCRRESTServer

	// Offer Manager
	OffersMgr *fcroffermgr.FCROfferMgr

	// Reputation Manager
	ReputationMgr *reputation.Reputation

	// RegistrationBlockHash is the hash of the block that registers this gateway
	// RegistrationTransactionReceipt is the transaction receipt containing the registration event
	// RegistrationMerkleRoot is the root of the merkle trie containing the transaction receipt
	// RegistrationMerkleProof proves the transaction receipt is part of the block
	RegistrationBlockHash          string
	RegistrationTransactionReceipt string
	RegistrationMerkleRoot         string
	RegistrationMerkleProof        *fcrmerkletree.FCRMerkleProof
}

// Single instance of the gateway
var instance *Core
var doOnce sync.Once

// GetSingleInstance returns the single instance of the gateway
func GetSingleInstance(confs ...*settings.AppSettings) *Core {
	doOnce.Do(func() {
		if len(confs) == 0 {
			logging.ErrorAndPanic("No settings supplied to Gateway start-up")
		}
		if len(confs) != 1 {
			logging.ErrorAndPanic("More than one sets of settings supplied to Gateway start-up")
		}

		instance = &Core{
			ProtocolVersion:                protocolVersion,
			ProtocolSupported:              []int32{protocolVersion, protocolSupported},
			Settings:                       confs[0],
			GatewayID:                      nil,
			GatewayPrivateKey:              nil,
			GatewayPrivateKeyVersion:       nil,
			OffersMgr:                      fcroffermgr.NewFCROfferMgr(),
			ReputationMgr:                  reputation.GetSingleInstance(),
			RegistrationBlockHash:          "TODO",
			RegistrationTransactionReceipt: "TODO",
			RegistrationMerkleRoot:         "TODO",
			RegistrationMerkleProof:        nil, //TODO
		}
	})
	return instance
}
