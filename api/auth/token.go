package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/wesleywcr/dev-book/api/config"
)

func CreateToken(userID uint64) (string, error) {
	permitions := jwt.MapClaims{}
	permitions["authorized"] = true
	permitions["exp"] = time.Now().Add(time.Hour * 6).Unix() //6hrs
	permitions["userId"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permitions)
	return token.SignedString([]byte(config.SecretKey))
}

func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, error := jwt.Parse(tokenString, returnVerificationKey)
	if error != nil {
		return error
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("Token inválido!")

}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func ExtractUserId(r *http.Request) (uint64, error) {
	tokenString := extractToken(r)
	token, error := jwt.Parse(tokenString, returnVerificationKey)
	if error != nil {
		return 0, error
	}
	if permitions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, error := strconv.ParseUint(fmt.Sprintf("%.0f", permitions["userId"]), 10, 64)
		if error != nil {
			return 0, error
		}
		return userId, nil

	}
	return 0, errors.New("Token inválido")
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
