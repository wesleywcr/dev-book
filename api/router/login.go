package router

import (
	"net/http"

	"github.com/wesleywcr/dev-book/api/controllers"
)

var routeLogin = Route{
	URI:                   "/login",
	Method:                http.MethodPost,
	HandleFunction:        controllers.Login,
	RequiredAuthorization: false,
}
