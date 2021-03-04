package reputation
// Copyright (C) 2020 ConsenSys Software Inc


const clientMaxReputation = int64(10000)
const clientMinReputaiton = int64(-10000)
const clientInitialReputation = int64(10)
const clientEstablishmentChallenge = int64(-1)
const clientOnChainDeposit = int64(1000)

//Response with one or more CID Offers. Initial payment and final payment made.
const clientStdDiscOneCidOffer = int64(10)

// Response with no CID Offers. Initial payment payment made.
const clientStdDiscNoCidOffers = int64(1)

// Response with one or more CID Offers. Response message sent after one second prior to TTL expiry. Initial payment payment made.
const clientStdDiscLateCidOffers = int64(1)

// Response with one or more CID Offers. Response message sent prior to one second prior to TTL expiry. Initial payment payment made but final payment not paid.
const clientStdDiscNonPayment = int64(-100)

// Response with one or more CID Offers from one or more Gateways. Initial payment and final payment made
const clientDhtDiscOneCidOffer = int64(10)

// Response with no CID Offers. Initial payment made.
const clientDhtDiscNoCidOffers = int64(1)

// Response with one or more CID Offers. Response message sent after one second prior to TTL expiry. Initial payment payment made.
const clientDhtDiscLateCidOffers = int64(1)

// Response with one or more CID Offers. Response message sent prior to one second prior to TTL expiry. Initial payment payment made but final payments not paid.
const clientDhtDiscNonPayment = int64(-300)

// Micro-payment paid for content via Gateway. Note that there will be many micro-payments during content retrieval.
const clientMicroPayment = int64(1)

// Invalid message received
const clientInvalidMessage = int64(-10)


