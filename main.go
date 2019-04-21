package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	ce "github.com/cloudevents/sdk-go"
)

func main() {

	port, err := strconv.Atoi(mustGetEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf("failed to parse port, %s", err.Error())
	}

	// Handler Mux
	mux := http.NewServeMux()

	// Ingres API Handler
	t, err := ce.NewHTTPTransport(
		ce.WithMethod("POST"),
		ce.WithPath("/"),
		ce.WithPort(port),
	)
	if err != nil {
		log.Fatalf("failed to create CloudEvents transport, %s", err.Error())
	}

	// wire handler for CE
	t.SetReceiver(&eventReceiver{})

	// Health Handler
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Events or UI Handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method, %s", r.Method)
		if r.Method == "POST" {
			t.ServeHTTP(w, r)
			return
		}
		fmt.Fprint(w, "Nothing to see here. Use POST to send CloudEvents")
	})

	a := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(a, mux))

}

func mustGetEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		log.Fatalf("Required env var (%s) not set", key)
	}

	return fallbackValue
}
