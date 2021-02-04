package nodeid

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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
)

const wordSize = 32 // 32 bytes

// NodeID represents a Gateway id
type NodeID struct {
	id []byte
}

// NewRandomNodeID creates a random node id object
func NewRandomNodeID() (*NodeID, error) {
	var n = NodeID{}
	n.id = make([]byte, wordSize)
	fcrcrypto.GeneratePublicRandomBytes(n.id)
	return &n, nil
}

// NewNodeID creates a node id object
func NewNodeID(id *big.Int) (*NodeID, error) {
	var n = NodeID{}
	b := id.Bytes()
	l := len(b)
	if l > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size1: %d", l)
	}
	idBytes := id.Bytes()
	n.id = make([]byte, wordSize)
	copy(n.id[wordSize-len(idBytes):], idBytes)
	return &n, nil
}

// NewNodeIDFromBytes creates a node id object
func NewNodeIDFromBytes(id []byte) (*NodeID, error) {
	var n = NodeID{}
	lenID := len(id)
	if lenID > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size2: %d", lenID)
	}
	n.id = make([]byte, wordSize)
	copy(n.id[wordSize-len(id):], id)
	return &n, nil
}

// NewNodeIDFromString creates a NodeID from a string
func NewNodeIDFromString(id string) (*NodeID, error) {
	var n = NodeID{}
	bytes, err := hex.DecodeString(id)
	if err != nil {
		return nil, err
	}

	if len(bytes) > wordSize {
		return nil, fmt.Errorf("NodeID: Incorrect size2: %d", len(id))
	}
	n.id = make([]byte, wordSize)
	copy(n.id[wordSize-len(bytes):], bytes)
	return &n, nil
}

// ToString returns a string for the node id.
func (n *NodeID) ToString() string {
	str := hex.EncodeToString(n.id)
	if str == "" {
		str = "00"
	}
	return str
}

// ToBytes returns the byte array representation of the node id.
func (n *NodeID) ToBytes() []byte {
	return n.id
}

// AsBytes32 returns the node id as a [32]byte
func (n *NodeID) AsBytes32() (result [wordSize]byte) {
	copy(result[:], n.id)
	return
}

// MarshalJSON is used to marshal NodeID into bytes
func (n NodeID) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.id)
}

// UnmarshalJSON is used to unmarshal bytes into NodeID
func (n *NodeID) UnmarshalJSON(p []byte) error {
	var id []byte
	err := json.Unmarshal(p, &id)
	if err != nil {
		return err
	}

	if len(id) != wordSize {
		return fmt.Errorf("NodeID: Incorrect size: %d", len(id))
	}
	n.id = make([]byte, wordSize)
	copy(n.id, id)
	return nil
}
