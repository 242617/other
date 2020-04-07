package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/242617/other/tasks/seabattle/sb"
)

type createMatrixRequest struct {
	Range int `json:"range"`
}

func (r *createMatrixRequest) Valid() bool {
	return r.Range > 0
}

func createMatrixHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request createMatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("cannot decode request: %s\n", err)
		http.Error(w, "cannot decode request", http.StatusBadRequest)
		return
	}

	if !request.Valid() {
		log.Println("invalid request")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var err error
	lock.Lock()
	game, err = sb.Create(request.Range)
	lock.Unlock()
	if err != nil {
		log.Printf("cannot create game: %s\n", err)
		http.Error(w, "cannot create game", http.StatusBadRequest)
		return
	}
}
