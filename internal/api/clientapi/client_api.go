package clientapi

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"log"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ant0ine/go-json-rest/rest"
)

const (
	clientAPIProtocolVersion = 1
	clientAPIProtocolSupportedHi = 1
)
// Can't have constant slices so create this at runtime.
// Order the API versions from most desirable to least desirable.
var clientAPIProtocolSupported []int


// ClientAPI holds the information for API between the Client and the Gateway.
type ClientAPI struct {
	// TODO: Add more fields (privkey, gateway id, etc.)
	// TODO: Add mutex for accessing gateway information.
}



// StartClientRestAPI starts the REST API as a separate go routine.
// Any start-up errors are returned.
func StartClientRestAPI(settings util.AppSettings) (*ClientAPI, error) {
	c := ClientAPI{}

	clientAPIProtocolSupported = make([]int, 1)
	clientAPIProtocolSupported[0] = clientAPIProtocolSupportedHi

	// Start the REST API and block until the error code is set.
	errChan := make(chan error, 1)
	go startRestAPI(settings, &c, errChan)
	return &c, <-errChan
}

func startRestAPI(settings util.AppSettings, c *ClientAPI, errChannel chan<- error) {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		// TODO: Remove these debug APIs prior to production release.
		rest.Get("/time", getTime),     // Get system time.
		rest.Get("/ip", getIP),         // Get IP address.
		rest.Get("/host", getHostname), // Get host name.

		rest.Post("/client/establishment", c.HandleClientNetworkEstablishment),       // Handle network establishment.
		rest.Post("/client/standard_request_cid", c.HandleClientStandardCIDDiscover), // Handle client standard cid request.
		rest.Post("/client/dht_request_cid", c.HandleClientDHTCIDDiscover),           // Handle DHT client cid request.
		rest.Post("/v1", c.msgRouter),
	)
	if err != nil {
		log.Fatal(err)
		errChannel <- err
		return
	}

	log.Println("Running REST API on: " + settings.BindRestAPI)
	api.SetApp(router)
	errChannel <- nil
	log.Fatal(http.ListenAndServe(":"+settings.BindRestAPI, api.MakeHandler()))
}


func (c *ClientAPI) msgRouter(w rest.ResponseWriter, r *rest.Request) {
	payload := messages.CommonRequestMessageFields{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		log.Println(err.Error())
		rest.Error(w, "Fail to decode payload.", http.StatusBadRequest)
		return
	}

	// Only process the rest of the message if the protocol version is understood.
	if (payload.ProtocolVersion != clientAPIProtocolVersion) {
		// Go through the protocol versions supported by the client and the
		// gateway from most desirable to least desirable, prioritising 
		// the gateway preference over the client preference.
		for _, clientProvVer := range payload.ProtocolSupported {
			for gatewayProtVer := range clientAPIProtocolSupported {
				if (clientProvVer == gatewayProtVer) {
					// Request the client switch to this protocol version
					response := messages.ProtocolChangeResponse{}
					response.MessageType = messages.ProtocolChange
					response.DesiredVersion = clientAPIProtocolVersion
					w.WriteJson(response)
					return
				}
			}
		}
		// No common protocol versions supported.
		response := messages.ProtocolMismatchResponse{}
		response.MessageType = messages.ProtocolMismatch
		w.WriteJson(response)
		return
	}

	switch (payload.MessageType) {
	case messages.ClientEstablishmentRequestType:
		c.HandleClientNetworkEstablishment(w, r)
	}



}