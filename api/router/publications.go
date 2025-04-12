package router

import (
	"net/http"

	"github.com/wesleywcr/dev-book/api/controllers"
)

var routesPublications = []Route{
	{
		URI:                   "/publications",
		Method:                http.MethodPost,
		HandleFunction:        controllers.CreatePublication,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/publications",
		Method:                http.MethodGet,
		HandleFunction:        controllers.GetPublications,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/publications/{publicationId}",
		Method:                http.MethodGet,
		HandleFunction:        controllers.GetPublicationsById,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/publications/{publicationId}",
		Method:                http.MethodPut,
		HandleFunction:        controllers.UpdatedPublication,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/publications/{publicationId}",
		Method:                http.MethodDelete,
		HandleFunction:        controllers.DeletePublication,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/users/{userId}/publications",
		Method:                http.MethodGet,
		HandleFunction:        controllers.SearchPublicationsByUserId,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/publications/{publicationId}/like",
		Method:                http.MethodPost,
		HandleFunction:        controllers.LikePublication,
		RequiredAuthorization: true,
	},
	{
		URI:                   "/publications/{publicationId}/deslike",
		Method:                http.MethodPost,
		HandleFunction:        controllers.DeslikePublication,
		RequiredAuthorization: true,
	},
}
