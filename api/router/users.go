package router

import (
	"net/http"

	"github.com/wesleywcr/dev-book/api/controllers"
)

var routesUsers = []Route{
	{
		URI:                   "/users",
		Method:                http.MethodPost,
		HandleFunction:        controllers.CreateUser,
		RequiredAuthorization: false,
	},
	{
		URI:                   "/users",
		Method:                http.MethodGet,
		HandleFunction:        controllers.ListUsers,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodGet,
		HandleFunction:        controllers.ListUser,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodPut,
		HandleFunction:        controllers.UpdateUser,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodDelete,
		HandleFunction:        controllers.DeleteUser,
		RequiredAuthorization: true,
	},
}
