package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func stateHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(game.State())
	if err != nil {
		log.Println("err", err)
	}
}
