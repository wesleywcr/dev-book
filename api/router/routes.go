package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wesleywcr/dev-book/api/middlewares"
)

type Route struct {
	URI                   string
	Method                string
	HandleFunction        func(http.ResponseWriter, *http.Request)
	RequiredAuthorization bool
}

func Config(r *mux.Router) *mux.Router {
	routes := routesUsers
	routes = append(routes, routeLogin)

	for _, route := range routes {
		if route.RequiredAuthorization {
			r.HandleFunc(route.URI, middlewares.Logger(middlewares.Authenticate(route.HandleFunction))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.HandleFunction)).Methods(route.Method)
		}
	}

	return r
}
