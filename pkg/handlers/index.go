package handlers

import (
	"log"
	"net/http"

	"github.com/mchmarny/myevents/pkg/utils"
)

// DefaultHandler handles index page
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	uid := getCurrentUserID(r)
	log.Printf("User ID: ", uid)

	// anon
	if uid == "" {
		if err := templates.ExecuteTemplate(w, "home", data); err != nil {
			log.Printf("Error in home template: %s", err)
		}
		return
	}

	// authenticated
	email, err := utils.ParseEmail(uid)
	if err != nil {
		log.Printf("Error parsing email: %v", err)
	}

	data["email"] = email
	if err := templates.ExecuteTemplate(w, "view", data); err != nil {
		log.Printf("Error in view template: %s", err)
	}

}
