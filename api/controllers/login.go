package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wesleywcr/dev-book/api/auth"
	"github.com/wesleywcr/dev-book/api/db"
	"github.com/wesleywcr/dev-book/api/models"
	"github.com/wesleywcr/dev-book/api/repositories"
	"github.com/wesleywcr/dev-book/api/response"
	"github.com/wesleywcr/dev-book/api/security"
)

// Login authenticates a user and returns a token.
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags Authentication
// @Accept json
// @Produce plain
// @Param credentials body models.User true "User credentials"
// @Success 200 {string} string "JWT Token"
// @Failure 422 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	bodyRequest, error := io.ReadAll(r.Body)
	if error != nil {
		response.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var user models.User
	if error = json.Unmarshal(bodyRequest, &user); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	repository := repositories.NewRepositoryOfUsers(db)
	userSalvedInDB, error := repository.SearchEmail(user.Email)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	if error = security.VerificatedPassoword(userSalvedInDB.Password, user.Password); error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}
	token, error := auth.CreateToken(userSalvedInDB.ID)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	w.Write([]byte(token))
}
