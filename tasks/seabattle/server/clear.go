package server

import "net/http"

func clearHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	game = nil
	lock.Unlock()
}
