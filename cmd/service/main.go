package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/mchmarny/myevents/pkg/handlers"
	"github.com/mchmarny/myevents/pkg/utils"
	"github.com/mchmarny/myevents/pkg/stores"
)

func main() {

	// Init db
	stores.InitDataStore()

	// Event Handlers
	http.HandleFunc("/", withLog(handlers.CloudEventHandler))

	// Server configured
	port := utils.MustGetEnv("PORT", "8080")

	log.Printf("Server starting on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func withLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(reqDump))
		}
		next.ServeHTTP(w, r)
	}
}
