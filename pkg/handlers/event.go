package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudevents/sdk-go/v02"
	"github.com/mchmarny/myevents/pkg/stores"
)

const (
	knownPublisherTokenName = "token"
)

// CloudEventHandler submitted messages
func CloudEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// check method
	if r.Method != http.MethodPost {
		log.Printf("wring method: %s", r.Method)
		http.Error(w, "Invalid method. Only POST supported", http.StatusMethodNotAllowed)
		return
	}

	// parse form to update
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, fmt.Sprintf("Post content error (%s)", err),
			http.StatusBadRequest)
		return
	}

	converter := v02.NewDefaultHTTPMarshaller()
	event, err := converter.FromRequest(r)
	if err != nil {
		log.Printf("error parsing cloudevent: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}
	log.Printf("Raw Event: %v", event)


	eventData, ok := event.Get("data")
	if !ok {
		http.Error(w, "Error, not a cloud event data", http.StatusBadRequest)
		return
	}
	log.Printf("Event Data: %v", eventData)

	saveErr := stores.SaveEvent(r.Context(), eventData)
	if saveErr != nil {
		log.Printf("error on event save: %v", saveErr)
		http.Error(w, "Error on event save", http.StatusBadRequest)
		return
	}

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(event)

}
