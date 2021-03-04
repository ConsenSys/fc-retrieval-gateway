package reputation

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// NOTE: At present all reputation is stored in memory. This will need to change as we move beyond a PoC.


func (r *Reputation) getClientReputation(clientNodeID *nodeid.NodeID) (val int64, exists bool) {
	clientNodeIDStr := clientNodeID.ToString()
	r.clientsMapLock.Lock()
	val, exists = r.clients[clientNodeIDStr]
	r.clientsMapLock.Unlock()
	return
}

func (r *Reputation) setClientReputation(clientNodeID *nodeid.NodeID, val int64) {
	clientNodeIDStr := clientNodeID.ToString()
	r.clientsMapLock.Lock()
	r.clients[clientNodeIDStr] = val
	r.clientsMapLock.Unlock()
}

func (r *Reputation) changeClientReputation(clientNodeID *nodeid.NodeID, amount int64) int64 {
	clientNodeIDStr := clientNodeID.ToString()
	var val int64
	var exists bool
	r.clientsMapLock.Lock()
	val, exists = r.clients[clientNodeIDStr]
	if (!exists) {
		r.clientsMapLock.Unlock()
		panic("changeClientReputation for non-existant client: " + clientNodeIDStr)
	}
	newVal := val + amount
	if (newVal > clientMaxReputation) {
		newVal = clientMaxReputation
	} else if (newVal < clientMinReputaiton) {
		newVal = clientMinReputaiton
	}

	r.clients[clientNodeIDStr] = newVal
	r.clientsMapLock.Unlock()

	return newVal
}
