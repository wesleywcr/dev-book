package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/wesleywcr/dev-book/api/auth"
	"github.com/wesleywcr/dev-book/api/db"
	"github.com/wesleywcr/dev-book/api/models"
	"github.com/wesleywcr/dev-book/api/repositories"
	"github.com/wesleywcr/dev-book/api/response"
	"github.com/wesleywcr/dev-book/api/security"
)

// CreateUser creates a new user.
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 422 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
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
	if error = user.Prepare("register"); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	//connect DB
	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	defer db.Close()
	// insert in DB
	repository := repositories.NewRepositoryOfUsers(db)
	user.ID, error = repository.Create(user)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusCreated, user)

}

// ListUsers retrieves all users.
// @Summary List all users
// @Description Retrieve a list of users filtered by name or nickname
// @Tags Users
// @Produce json
// @Param user query string false "Name or nickname to filter"
// @Success 200 {array} models.User
// @Failure 500 {object} response.ErrorResponse
// @Router /users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNickname := strings.ToLower(r.URL.Query().Get("user"))

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	respository := repositories.NewRepositoryOfUsers(db)

	users, error := respository.Search(nameOrNickname)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// ListUser retrieves a specific user by ID.
// @Summary Get a user by ID
// @Description Retrieve a user's details by their ID
// @Tags Users
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId} [get]
func ListUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)

	user, error := repository.SearchPerId(userId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusOK, user)

}

// UpdateUser updates a user's information.
// @Summary Update a user
// @Description Update the details of an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	userIdToken, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}
	if userId != userIdToken {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não é o seu"))
		return
	}

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

	if error := user.Prepare("update"); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)
	if error = repository.Update(userId, user); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)

}

// DeleteUser deletes a user.
// @Summary Delete a user
// @Description Remove a user from the system
// @Tags Users
// @Produce json
// @Param userId path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	userIdToken, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}
	if userId != userIdToken {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não é o seu"))
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)
	if error = repository.Delete(userId); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)

}

// FollowUser allows a user to follow another user.
// @Summary Follow a user
// @Description Follow another user by their ID
// @Tags Users
// @Produce json
// @Param userId path int true "User ID to follow"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/follow [post]
func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}
	if followerId == userId {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo"))
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)
	if error = repository.Follow(userId, followerId); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)

}

// UnFollowUser allows a user to unfollow another user.
// @Summary Unfollow a user
// @Description Stop following a user by their ID
// @Tags Users
// @Produce json
// @Param userId path int true "User ID to unfollow"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/unfollow [post]
func UnFollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}
	if followerId == userId {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível deixar de seguir você mesmo"))
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)
	if error = repository.UnFollow(userId, followerId); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)

}

// SearchFollowers retrieves the followers of a user.
// @Summary Get followers
// @Description Retrieve a list of users following a specific user
// @Tags Users
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} models.User
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/followers [get]
func SearchFollowers(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)
	followers, error := repository.SearchFollowers(userId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusOK, followers)

}

// SearchFollowing retrieves the users a specific user is following.
// @Summary Get following users
// @Description Retrieve a list of users a specific user is following
// @Tags Users
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} models.User
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/following [get]
func SearchFollowing(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)

	users, error := repository.SearchFollowing(userId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusOK, users)

}

// UpdatePassword updates a user's password.
// @Summary Update password
// @Description Change the password of a user
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param password body models.Password true "Password data"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/password [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	userIdToken, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}

	parameters := mux.Vars(r)
	userId, error := strconv.ParseUint(parameters["userId"], 10, 64)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}
	if userIdToken != userId {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível atualizar senha de um usuário que não seja o seu"))
		return
	}

	bodyRequest, error := io.ReadAll(r.Body)

	var password models.Password

	if error = json.Unmarshal(bodyRequest, &password); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfUsers(db)

	passwordSavedDB, error := repository.GetPassword(userId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	if error = security.VerificatedPassoword(passwordSavedDB, password.Current); error != nil {
		response.Error(w, http.StatusInternalServerError, errors.New("A senha atual está incorreta"))
		return
	}

	passwordWithHash, error := security.Hash(password.New)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	if error := repository.UpdatePassword(userId, string(passwordWithHash)); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}
