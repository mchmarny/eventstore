package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mchmarny/myevents/pkg/handlers"
	"github.com/mchmarny/myevents/pkg/utils"
	"github.com/mchmarny/myevents/pkg/stores"
	"github.com/knative/pkg/cloudevents"
)

func main() {

	// Init db
	stores.InitDataStore()

	// Event Handlers
	m := cloudevents.NewMux()
	err := m.Handle("tech.knative.demo.kapi.stock", handlers.CloudEventHandler)
	if err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	err = m.Handle("io.redis.queue", handlers.CloudEventHandler)
	if err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	// Server configured
	port := utils.MustGetEnv("PORT", "8080")

	log.Printf("Server starting on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), m))

}
