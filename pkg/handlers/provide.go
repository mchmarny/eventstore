package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	eventChannel = make(chan []byte, 1)
)

// WSHandler provides backing service for the UI
func WSHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("WS connection...")

	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Printf("Error while upgrading WS: %v", err)
		http.Error(w, "Error while upgrading WS", http.StatusBadRequest)
		return
	}

	for {
		select {
		case m := <-eventChannel:
			if connErr := conn.WriteMessage(websocket.TextMessage, m); connErr != nil {
				log.Printf("Error on write message: %v", connErr)
			}
		}
	}

}
