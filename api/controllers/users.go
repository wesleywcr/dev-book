package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/wesleywcr/dev-book/api/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, error := io.ReadAll(r.Body)
	if error != nil {
		log.Fatal(error)
	}
	var user models.User
	if error = json.Unmarshal(bodyRequest, &user); error != nil {
		log.Fatal((error))
	}

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
