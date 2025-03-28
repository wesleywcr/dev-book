package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryOfPublications(db)
	publication, error := repository.SearchPublicationsById(publicationId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusOK, publication)

}
func SearchPublication(w http.ResponseWriter, r *http.Request) {
	userId, error := auth.ExtractUserId(r)
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

	repository := repositories.NewRepositoryOfPublications(db)
	publications, error := repository.Search(userId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusOK, publications)
}

func UpdatedPublication(w http.ResponseWriter, r *http.Request) {
	userId, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryOfPublications(db)

	publicationSalvedDB, error := repository.SearchPublicationsById(publicationId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	if publicationSalvedDB.AuthorID != userId {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível atualizar uma publicação que não seja a sua"))
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

	if error = publication.Prepare(); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	if error = repository.Update(publicationId, publication); error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeletePublication(w http.ResponseWriter, r *http.Request) {
	userId, error := auth.ExtractUserId(r)
	if error != nil {
		response.Error(w, http.StatusBadRequest, error)
		return
	}

	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryOfPublications(db)

	publicationSalvedDB, error := repository.SearchPublicationsById(publicationId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	if publicationSalvedDB.AuthorID != userId {
		response.Error(w, http.StatusForbidden, errors.New("Não é possível deletar uma publicação que não seja a sua"))
		return
	}
	if error := repository.Delete(publicationId); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func SearchPublicationsByUserId(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewRepositoryOfPublications(db)
	publications, error := repository.SearchPublicationByUserId(userId)
	if error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}
	response.JSON(w, http.StatusOK, publications)
}

func LikePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryOfPublications(db)

	if error := repository.Like(publicationId); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
func DeslikePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationId, error := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryOfPublications(db)

	if error := repository.Deslike(publicationId); error != nil {
		response.Error(w, http.StatusInternalServerError, error)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
