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
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	NewRandomNodeID()
}

func TestRoundTripBigInt(t *testing.T) {
	testRoundTripBigInt(t, "00")
	testRoundTripBigInt(t, "01")
	testRoundTripBigInt(t, "fe")
	testRoundTripBigInt(t, "0100")
	testRoundTripBigInt(t, "30010203040506070809")
	testRoundTripBigInt(t, "80010203040506070809")
}

func TestRoundTripFromBytes(t *testing.T) {
	testRoundTripFromBytes(t, "00")
	testRoundTripFromBytes(t, "01")
	testRoundTripFromBytes(t, "fe")
	testRoundTripFromBytes(t, "0100")
	testRoundTripFromBytes(t, "30010203040506070809")
	testRoundTripFromBytes(t, "80010203040506070809")
}

func TestRoundTripFromString(t *testing.T) {
	testRoundTripFromString(t, "00")
	testRoundTripFromString(t, "01")
	testRoundTripFromString(t, "fe")
	testRoundTripFromString(t, "0100")
	testRoundTripFromString(t, "30010203040506070809")
	testRoundTripFromString(t, "80010203040506070809")
}

func testRoundTripBigInt(t *testing.T, value string) {
	id := new(big.Int)
	_, ok := id.SetString(value, 16)
	if !ok {
		panic("Number format issue: " + value)
	}

	nodeID, err := NewNodeID(id)
	if err != nil {
		panic(err)
	}
	idStr := nodeID.ToString()
	assert.Equal(t, fmt.Sprintf("%064s", value), idStr, "NewNodeID failed")

}

func testRoundTripFromBytes(t *testing.T, value string) {
	bytes, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}

	nodeID, err := NewNodeIDFromBytes(bytes)
	if err != nil {
		panic(err)
	}
	idStr := nodeID.ToString()
	assert.Equal(t, fmt.Sprintf("%064s", value), idStr, "NewNodeIDFromString failed")

}

func testRoundTripFromString(t *testing.T, value string) {
	nodeID, err := NewNodeIDFromString(value)
	if err != nil {
		panic(err)
	}
	idStr := nodeID.ToString()
	assert.Equal(t, fmt.Sprintf("%064s", value), idStr, "NewNodeIDFromString failed")

}
