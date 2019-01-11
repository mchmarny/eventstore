package handlers

import (
	"log"
	"net/http"
)

// DefaultHandler handles index page
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}

	uid := getCurrentUserID(r)
	// TODO: Redirect to queue view page if already authenticated
	log.Printf("User authenticated: %s, getting data...", uid)

	if err := templates.ExecuteTemplate(w, "home", data); err != nil {
		log.Printf("Error in home template: %s", err)
	}

}
