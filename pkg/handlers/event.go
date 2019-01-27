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

	m := v02.NewDefaultHTTPMarshaller()
	e, err := m.FromRequest(r)
	if err != nil {
		log.Printf("error parsing cloudevent: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}

	log.Printf("Raw Event: %v", e)
	data, ok := e.Get("data")
	if !ok {
		log.Println("nil event data [data]")
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}

	log.Printf("Inner event v0.2: %v", data)
	e2 := data.(map[string]interface{})
	e2ID := e2["id"]
	if e2ID == nil {
		log.Println("nil event ID [id]")
		http.Error(w, fmt.Sprintf("Invalid Cloud Event (%v)", err),
			http.StatusBadRequest)
		return
	}

	saveErr := stores.SaveEvent(r.Context(), e2ID.(string), e2)
	if saveErr != nil {
		log.Printf("error on event save: %v", saveErr)
		http.Error(w, "Error on event save", http.StatusBadRequest)
		return
	}

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(e2)

}
