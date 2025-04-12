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

// CreatePublication creates a new publication.
// @Summary Create a publication
// @Description Create a new publication for the authenticated user
// @Tags Publications
// @Accept json
// @Produce json
// @Param publication body models.Publication true "Publication data"
// @Success 201 {object} models.Publication
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications [post]
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

// GetPublications retrieves a publication by ID.
// @Summary Get a publication
// @Description Retrieve a publication by its ID
// @Tags Publications
// @Produce json
// @Param publicationId path int true "Publication ID"
// @Success 200 {object} models.Publication
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications/{publicationId} [get]
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

// SearchPublication retrieves all publications for the authenticated user.
// @Summary List publications
// @Description Retrieve all publications for the authenticated user
// @Tags Publications
// @Produce json
// @Success 200 {array} models.Publication
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications [get]
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

// UpdatedPublication updates a publication.
// @Summary Update a publication
// @Description Update a publication owned by the authenticated user
// @Tags Publications
// @Accept json
// @Produce json
// @Param publicationId path int true "Publication ID"
// @Param publication body models.Publication true "Updated publication data"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications/{publicationId} [put]
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

// DeletePublication deletes a publication.
// @Summary Delete a publication
// @Description Delete a publication owned by the authenticated user
// @Tags Publications
// @Produce json
// @Param publicationId path int true "Publication ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications/{publicationId} [delete]
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

// SearchPublicationsByUserId retrieves publications by a specific user.
// @Summary Get publications by user
// @Description Retrieve all publications created by a specific user
// @Tags Publications
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} models.Publication
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userId}/publications [get]
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

// LikePublication likes a publication.
// @Summary Like a publication
// @Description Add a like to a publication
// @Tags Publications
// @Produce json
// @Param publicationId path int true "Publication ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications/{publicationId}/like [post]
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

// DeslikePublication removes a like from a publication.
// @Summary Unlike a publication
// @Description Remove a like from a publication
// @Tags Publications
// @Produce json
// @Param publicationId path int true "Publication ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /publications/{publicationId}/unlike [post]
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
