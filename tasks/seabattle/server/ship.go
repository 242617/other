package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/242617/other/tasks/seabattle/sb"
)

type shipRequest struct {
	Coordinates string `json:"Coordinates"`
}

func shipHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request shipRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("cannot decode request: %s\n", err)
		http.Error(w, "cannot decode request", http.StatusBadRequest)
		return
	}

	var ships []*sb.Ship

	pairs := strings.Split(request.Coordinates, ",")
	for _, pair := range pairs {

		rawPoints := strings.Split(pair, " ")
		if len(rawPoints) != 2 {
			log.Printf(`cannot parse pair "%s"\n`, pair)
			http.Error(w, "cannot parse pair", http.StatusBadRequest)
			return
		}

		begin, err := sb.ParsePoint(rawPoints[0])
		if err != nil {
			log.Printf(`cannot parse begin point "%s": "%s"\n`, rawPoints[0], err)
			http.Error(w, "cannot parse begin point", http.StatusBadRequest)
			return
		}

		end, err := sb.ParsePoint(rawPoints[1])
		if err != nil {
			log.Printf(`cannot parse end point "%s": "%s"\n`, rawPoints[1], err)
			http.Error(w, "cannot parse end point", http.StatusBadRequest)
			return
		}

		ships = append(ships, sb.NewShip(*begin, *end))

	}

	lock.Lock()
	defer lock.Unlock()
	if err := game.Setup(ships); err != nil {
		log.Printf("cannot setup game: %s", err)
		http.Error(w, "cannot setup game", http.StatusInternalServerError)
		return
	}
}
