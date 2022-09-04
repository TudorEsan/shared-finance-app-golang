package sharedmodels

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY []byte = GetSecretKey()

type SignedDetails struct {
	Email    string
	Username string
	Id       string
	jwt.StandardClaims
}

func GetSecretKey() []byte {
	key, ok := os.LookupEnv("SECRET_JWT")
	if !ok {
		log.Println("SECRET KEY WAS INITIALIZED BY DEFAULT")
		return []byte("TEST")
	}
	return []byte(key)
}
