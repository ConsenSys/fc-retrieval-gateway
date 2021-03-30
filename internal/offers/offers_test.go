package offers

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
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestOffersInitial(t *testing.T) {
	GetSingleInstance()
}

func TestEmpty(t *testing.T) {
	o := GetSingleInstance()

	_, exists := o.GetOffers(cidOne(t))
	assert.False(t, exists, "Found CID offer despite empty offer system")
}

func TestAddOld(t *testing.T) {
	o := GetSingleInstance()
	err := o.Add(createOldSingleCidGroupOfferCidOne(t))
	if err == nil {
		t.Errorf("Didn't return error due to CID having expired.")
	}

}

func TestAdd(t *testing.T) {
	o := newInstance()
	err := o.Add(createNewSingleCidGroupOfferCidOne(t))
	if err != nil {
		t.Errorf("Error returned by Add: %e", err)
	}

	cidOffers, exists := o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")

	o.ExpireOffers()
	cidOffers, exists = o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")
}

func TestAddTwo(t *testing.T) {
	o := newInstance()
	o.Add(createNewSingleCidGroupOfferCidOne(t))
	o.Add(createFutureSingleCidGroupOfferCidOne(t))

	cidOffers, exists := o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 2, len(cidOffers), "Should be only two CID offers")
}

func TestExpire(t *testing.T) {
	o := newInstance()
	o.Add(createNewSingleCidGroupOfferCidOne(t))

	_, exists := o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID")

	now := time.Now()
	nowSeconds := now.Unix()
	mockNow := nowSeconds + 1001
	util.SetMockedClock(mockNow)

	o.ExpireOffers()
	util.SetRealClock()

	_, exists = o.GetOffers(cidOne(t))
	assert.False(t, exists, "Found offers for CID when it should have expired")
}

func TestExpireTwo(t *testing.T) {
	o := newInstance()
	o.Add(createNewSingleCidGroupOfferCidOne(t))
	o.Add(createFutureSingleCidGroupOfferCidOne(t))

	cidOffers, exists := o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID")
	assert.Equal(t, 2, len(cidOffers), "Should be only two CID offers")

	now := time.Now()
	nowSeconds := now.Unix()
	mockNow := nowSeconds + 1001
	util.SetMockedClock(mockNow)

	o.ExpireOffers()

	cidOffers, exists = o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")
	cidExp := cidOffers[0].GetExpiry()
	assert.Greater(t, cidExp, mockNow, "Incorrect CID offer expired")

	mockNow = nowSeconds + 2001
	util.SetMockedClock(mockNow)

	o.ExpireOffers()
	util.SetRealClock()

	cidOffers, exists = o.GetOffers(cidOne(t))
	assert.False(t, exists, "Found offers for CID when they should have expired")
}

func TestMultipleCids(t *testing.T) {
	o := newInstance()
	err := o.Add(createNewCidGroupOfferCidMultiple(t))
	if err != nil {
		t.Errorf("Error returned by Add: %e", err)
	}

	cidOffers, exists := o.GetOffers(cidOne(t))
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")

	cidOffers, exists = o.GetOffers(cidTwo(t))
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")

	cidOffers, exists = o.GetOffers(cidThree(t))
	assert.True(t, exists, "Can't find any offers for CID 1")
	assert.Equal(t, 1, len(cidOffers), "Should be only one CID offer")
}

func createOldSingleCidGroupOfferCidOne(t *testing.T) *cidoffer.CIDOffer {
	return createSingleCidGroupOffer(t, cidOne(t), 0)
}

func createNewSingleCidGroupOfferCidOne(t *testing.T) *cidoffer.CIDOffer {
	return createSingleCidGroupOffer(t, cidOne(t), 1)
}

func createFutureSingleCidGroupOfferCidOne(t *testing.T) *cidoffer.CIDOffer {
	return createSingleCidGroupOffer(t, cidOne(t), 2)
}

func createNewCidGroupOfferCidMultiple(t *testing.T) *cidoffer.CIDOffer {
	cids := make([]cid.ContentID, 0)
	cids = append(cids, *cidOne(t))
	cids = append(cids, *cidTwo(t))
	cids = append(cids, *cidThree(t))
	return createCidGroupOffer(t, cids, 1)
}

func createSingleCidGroupOffer(t *testing.T, theCid *cid.ContentID, howNew int) *cidoffer.CIDOffer {
	cids := make([]cid.ContentID, 0)
	cids = append(cids, *theCid)
	return createCidGroupOffer(t, cids, howNew)
}

func createCidGroupOffer(t *testing.T, cids []cid.ContentID, howNew int) *cidoffer.CIDOffer {
	aNodeID := nodeid.NewRandomNodeID()
	price := uint64(5)
	now := time.Now()
	nowSeconds := now.Unix()
	var expiry int64
	switch howNew {
	case 0:
		expiry = nowSeconds - 1
	case 2:
		expiry = nowSeconds + 2000
	default:
		expiry = nowSeconds + 1000
	}
	qos := uint64(5)
	c, err := cidoffer.NewCIDOffer(aNodeID, cids, price, expiry, qos)
	if err != nil {
		t.Errorf("Error returned by NewCidGroupOffer: %e", err)
	}
	return c
}

func cidOne(t *testing.T) *cid.ContentID {
	cid, err := cid.NewContentID(big.NewInt(1))
	if err != nil {
		t.Errorf("Error returned by NewContentID for CID %cid: %e", cid, err)
	}
	return cid
}

func cidTwo(t *testing.T) *cid.ContentID {
	cid, err := cid.NewContentID(big.NewInt(2))
	if err != nil {
		t.Errorf("Error returned by NewContentID for CID %cid: %e", cid, err)
	}
	return cid
}

func cidThree(t *testing.T) *cid.ContentID {
	cid, err := cid.NewContentID(big.NewInt(3))
	if err != nil {
		t.Errorf("Error returned by NewContentID for CID %cid: %e", cid, err)
	}
	return cid
}
