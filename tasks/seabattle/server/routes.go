package server

import (
	"net/http"

	"github.com/242617/other/tasks/seabattle/sm"
)

type route struct {
	Method  string
	Action  sm.Action
	Handler func(w http.ResponseWriter, r *http.Request)
}

type request struct {
	url    string
	method string
}
type handler struct {
	action  sm.Action
	handler func(w http.ResponseWriter, r *http.Request)
}
type routes map[request]handler

func (r routes) Find(url, method string) (handler, bool) {
	for requestKey, handlerValue := range r {
		if requestKey.url == url && requestKey.method == method {
			return handlerValue, true
		}
	}
	return handler{}, false
}
