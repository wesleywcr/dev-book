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
)

func CreatePublication(w http.ResponseWriter, r *http.Request) {
	userId, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusUnauthorized, error)
		return
	}
	bodyRequest, error := io.ReadAll(r.Body)
	if error != nil {
		response.Error(w, http.StatusUnprocessableEntity, error)
		return
	}

	var publication models.Publication

	if error = json.Unmarshal(bodyRequest, &publication); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	publication.AuthorID = userId
	if error := publication.Prepare(); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}
	db, error := db.ConnectDB()
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryOfPublications(db)
	publication.ID, error = repository.Create(publication)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusCreated, publication)
}
func GetPublications(w http.ResponseWriter, r *http.Request) {

}
func SearchPublication(w http.ResponseWriter, r *http.Request) {

}
func Updateublication(w http.ResponseWriter, r *http.Request) {

}
func DeletePublication(w http.ResponseWriter, r *http.Request) {

}
