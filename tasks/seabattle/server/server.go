package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/242617/other/tasks/seabattle/sb"
	"github.com/242617/other/tasks/seabattle/sm"
)

var Routes = routes{
	{"/create-matrix", http.MethodPost}: {sm.ActionCreateMatrix, createMatrixHandler},
	{"/ship", http.MethodPost}:          {sm.ActionShip, shipHandler},
	{"/shot", http.MethodPost}:          {sm.ActionShot, shotHandler},
	{"/clear", http.MethodPost}:         {sm.ActionClear, clearHandler},
	{"/state", http.MethodGet}:          {sm.ActionState, stateHandler},
}

var (
	game *sb.SeaBattle
	lock sync.Mutex
)
var machine *sm.StateMachine

func Init(m *sm.StateMachine) error {
	machine = m
	return nil
}

func Start(addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler, ok := Routes.Find(r.RequestURI, r.Method)
		if !ok {
			log.Printf("route not found: url: %s, method: %s\n", r.RequestURI, r.Method)
			http.Error(w, "incorrect route", http.StatusBadRequest)
			return
		}
		if ok := machine.Action(handler.action); !ok {
			log.Printf("cannot change state: state %s, action: %s\n", machine.Current(), handler.action)
			http.Error(w, "incorrect request", http.StatusBadRequest)
			return
		}
		handler.handler(w, r)
	})
	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		state := struct {
			MachineState interface{}
			GameState    interface{}
		}{}
		state.MachineState = machine.Current()
		if game != nil {
			state.GameState = game.State()
		}
		json.NewEncoder(w).Encode(state)
	})
	return http.ListenAndServe(addr, nil)
}
