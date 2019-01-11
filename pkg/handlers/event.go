package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mchmarny/myvents/pkg/utils"
	"github.com/knative/pkg/cloudevents"
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

	// check for presense of publisher token
	srcToken := r.URL.Query().Get(knownPublisherTokenName)
	if srcToken == "" {
		log.Printf("nil token: %s", srcToken)
		http.Error(w, fmt.Sprintf("Invalid request (%s missing)", knownPublisherTokenName),
			http.StatusBadRequest)
		return
	}

	// check validity of poster token
	if !utils.Contains(knownPublisherTokens, srcToken) {
		log.Printf("invalid token: %s", srcToken)
		http.Error(w, fmt.Sprintf("Invalid publisher token value (%s)", knownPublisherTokenName),
			http.StatusBadRequest)
		return
	}

	var eventData string
	eventCtx, err := cloudevents.FromRequest(eventData, r)
	if err != nil {
		log.Printf("error parsing cloudevent: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}

	// TODO: push the event to UI display channel
	log.Printf("Event Context: %+v", eventCtx)
	log.Printf("Event Data: %s", eventData)

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(eventData)

}
