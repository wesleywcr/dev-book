package auth

import (
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
