package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI                   string
	Method                string
	HandleFunction        func(http.ResponseWriter, *http.Request)
	RequiredAuthorization bool
}

func Config(r *mux.Router) *mux.Router {
	routes := routesUsers
	for _, route := range routes {
		r.HandleFunc(route.URI, route.HandleFunction).Methods(route.Method)
	}
	return r
}
