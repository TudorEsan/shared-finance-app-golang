package sharedhelpers

import (
	"os"
	"strings"

	"github.com/TudorEsan/shared-finance-app-golang/customErrors"
	sharedmodels "github.com/TudorEsan/shared-finance-app-golang/sharedModels"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(signedToken string) (*sharedmodels.SignedDetails, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET not set")
	}
	token, err := jwt.ParseWithClaims(signedToken, &sharedmodels.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil && strings.Contains(err.Error(), "expired") {
		return nil, customErrors.ExpiredToken{}
	}
	if err != nil {
		return nil, customErrors.InvalidToken{E: err}
	}

	claims, ok := token.Claims.(*sharedmodels.SignedDetails)
	if !ok {
		return nil, customErrors.InvalidToken{}
	}
	return claims, nil
}

func RemoveCookies(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", "", 60*60*24*30, "", "", false, false)
}
