package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mchmarny/kueue/handlers"
	"github.com/mchmarny/kueue/utils"
)


func main() {

	// init configs
	handlers.InitHandlers()


	// Mux
	mux := http.NewServeMux()

	// Static
	mux.Handle("/static/", http.StripPrefix("/static/",
		  http.FileServer(http.Dir("static"))))

	// Handlers
	mux.HandleFunc("/", handlers.DefaultHandler)
	mux.HandleFunc("/auth/login", handlers.OAuthLoginHandler)
	mux.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler)
	mux.HandleFunc("/auth/logout", handlers.LogOutHandler)
	mux.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Server configured
	port := utils.MustGetEnv("PORT", "8080")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Printf("Server starting on port %s \n", port)
	log.Fatal(server.ListenAndServe())

}
