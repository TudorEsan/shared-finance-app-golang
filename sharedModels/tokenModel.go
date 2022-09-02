package sharedmodels

import (
	"os"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY string = GetSecretKey()

type SignedDetails struct {
	Email    string
	Username string
	Id       string
	jwt.StandardClaims
}

func GetSecretKey() string {
	return os.Getenv("SECRET_JWT")
}