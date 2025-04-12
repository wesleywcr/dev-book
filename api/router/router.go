package router

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/wesleywcr/dev-book/api/docs" // Import generated Swagger docs
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return Config(r)
}
