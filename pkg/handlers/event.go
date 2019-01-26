package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ce "github.com/knative/pkg/cloudevents"
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

	var event map[string]string
	ctx, err := ce.FromRequest(event, r)
	if err != nil {
		log.Printf("error parsing cloudevent: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}
	log.Printf("Event Context: %v", ctx)
	log.Printf("Raw Event Data: %v", event)

	ctx.Extensions = map[string]interface{}{ "raw": event }

	saveErr := stores.SaveEvent(r.Context(), ctx)
	if saveErr != nil {
		log.Printf("error on event save: %v", saveErr)
		http.Error(w, "Error on event save", http.StatusBadRequest)
		return
	}

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(event)

}
