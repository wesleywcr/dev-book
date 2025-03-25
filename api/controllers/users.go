package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wesleywcr/dev-book/api/db"
	"github.com/wesleywcr/dev-book/api/models"
	"github.com/wesleywcr/dev-book/api/repositories"
	"github.com/wesleywcr/dev-book/api/response"
)

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
	if error = user.Prepare(); error != nil {
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
func ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get users!"))
}

func ListUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get user!"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update users!"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deletar users!"))
}
