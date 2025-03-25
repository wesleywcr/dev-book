package router

import (
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {

	r := mux.NewRouter()
	return Config(r)
}
