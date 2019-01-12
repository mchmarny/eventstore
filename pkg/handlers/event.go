package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudevents/sdk-go/v02"
	"github.com/mchmarny/myevents/pkg/utils"
	"github.com/mchmarny/myevents/pkg/queue"
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

	converter := v02.NewDefaultHTTPMarshaller()
	event, err := converter.FromRequest(r)
	if err != nil {
		log.Printf("error parsing cloudevent: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}

	log.Printf("Event: %v", event)


	// content from the event
	eventContent, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error while marshaling event: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}

	// publish event
	if !queue.GetQueue().PublishBytes(eventContent) {
		log.Println("error publishing event")
		http.Error(w, "error publishing event", http.StatusBadRequest)
		return
	}

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(event)

}
