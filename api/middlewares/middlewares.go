package middlewares

import (
	"log"
	"net/http"

	"github.com/wesleywcr/dev-book/api/auth"
	"github.com/wesleywcr/dev-book/api/response"
)

func Logger(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		nextFunction(w, r)
	}
}
func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if error := auth.ValidateToken(r); error != nil {
			response.Error(w, http.StatusUnauthorized, error)
			return
		}
		nextFunction(w, r)
	}
}
