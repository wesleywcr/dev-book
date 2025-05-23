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
	{
		URI:                   "/users/{userId}/follow",
		Method:                http.MethodPost,
		HandleFunction:        controllers.FollowUser,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}/unfollow",
		Method:                http.MethodPost,
		HandleFunction:        controllers.UnFollowUser,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}/followers",
		Method:                http.MethodGet,
		HandleFunction:        controllers.SearchFollowers,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}/following",
		Method:                http.MethodGet,
		HandleFunction:        controllers.SearchFollowing,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}/update-password",
		Method:                http.MethodPost,
		HandleFunction:        controllers.UpdatePassword,
		RequiredAuthorization: true,
	},
}
