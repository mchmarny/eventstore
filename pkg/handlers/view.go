package handlers

import (
	"log"
	"net/http"

	"github.com/mchmarny/myevents/pkg/utils"
)

// ViewHandler handles view page
func ViewHandler(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	uid := getCurrentUserID(r)
	log.Printf("User ID: %s", uid)

	// anon
	if uid == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
