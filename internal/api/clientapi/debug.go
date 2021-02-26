package clientapi

// Copyright (C) 2020 ConsenSys Software Inc

// Contains debug APIs

import (
	"net"
	"os"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ant0ine/go-json-rest/rest"
)

func getTime(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(time.Now())
}

func getHostname(w rest.ResponseWriter, r *rest.Request) {
	name, err := os.Hostname()
	if err != nil {
		logging.Info("Get host name1: %s", err.Error())
		return
	}

	w.WriteJson(name)
}

func getIP(w rest.ResponseWriter, r *rest.Request) {
	name, err := os.Hostname()
	if err != nil {
		logging.Info("Get host name2: %s", err.Error())
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		logging.Info("Lookup host: %s", err.Error())
		return
	}

	w.WriteJson(addrs)
}
