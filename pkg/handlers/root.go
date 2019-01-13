package handlers

import (
	"log"
	"net/http"
)

// RootHandler handles view page
func RootHandler(w http.ResponseWriter, r *http.Request) {

	// if POST on root
	if r.Method == http.MethodPost {
		CloudEventHandler(w, r)
		return
	}

	data := make(map[string]interface{})

	uid := getCurrentUserID(r)
	log.Printf("User ID: %s", uid)

	// authenticated
	if uid != "" {
		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return
	}

	// anonymous
	if err := templates.ExecuteTemplate(w, "home", data); err != nil {
		log.Printf("Error in home template: %s", err)
	}

}
