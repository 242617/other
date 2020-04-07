package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/242617/other/tasks/seabattle/sb"
)

type shotRequest struct {
	Coordinate string `json:"coord"`
}

func shotHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request shotRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("cannot decode request: %s\n", err)
		http.Error(w, "cannot decode request", http.StatusBadRequest)
		return
	}

	point, err := sb.ParsePoint(request.Coordinate)
	if err != nil {
		log.Printf("cannot parse point: %s\n", err)
		http.Error(w, "cannot parse point", http.StatusBadRequest)
		return
	}

	lock.Lock()
	info, err := game.Shot(*point)
	lock.Unlock()
	if err != nil {
		log.Printf("cannot make shot: %s\n", err)
		http.Error(w, "cannot make shot", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(info)
	if err != nil {
		log.Println("err", err)
		return
	}
}
