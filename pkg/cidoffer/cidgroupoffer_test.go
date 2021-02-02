package cidoffer

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
	"math/big"
	"testing"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

func TestGetPrice(t *testing.T) {
	aNodeID, err := nodeid.NewNodeID(big.NewInt(7))
	if err != nil {
		panic(err)
	}
	aCid, ciderr := cid.NewContentID(big.NewInt(7))
	cids := make([]cid.ContentID, 0)
	cids = append(cids, *aCid)
	price := uint64(5)
	expiry := int64(10)
	qos := uint64(5)
	c, cidgerr := NewCidGroupOffer(aNodeID, &cids, price, expiry, qos)
	if ciderr != nil {
		t.Errorf("Error returned by NewContentID: %e", err)
	}
	if cidgerr != nil {
		t.Errorf("Error returned by NewCidGroupOffer: %e", err)
	}
	if c.GetPrice() != price {
		t.Errorf("Expected: %d, Actual: %d", price, c.GetPrice())
	}
}
