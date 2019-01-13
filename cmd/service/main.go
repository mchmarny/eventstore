package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mchmarny/myevents/pkg/handlers"
	"github.com/mchmarny/myevents/pkg/utils"

	"golang.org/x/net/websocket"
)

func main() {

	// init configs
	handlers.InitHandlers()


	// Static
	http.Handle("/static/", http.StripPrefix("/static/",
		  http.FileServer(http.Dir("static"))))

	// UI Handlers
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/view", handlers.ViewHandler)
	http.Handle("/ws", websocket.Handler(handlers.WSHandler))

	// Auth Handlers
	http.HandleFunc("/auth/login", handlers.OAuthLoginHandler)
	http.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler)
	http.HandleFunc("/auth/logout", handlers.OAuthLogoutHandler)

	// Ingres API Handler
	http.HandleFunc("/v1/event", handlers.CloudEventHandler)

	// Health Handler
	http.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Server configured
	port := utils.MustGetEnv("PORT", "8080")

	log.Printf("Server starting on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}
