package sharedvalidators

import (
	"strings"

	"github.com/TudorEsan/shared-finance-app-golang/customErrors"
	sharedmodels "github.com/TudorEsan/shared-finance-app-golang/sharedModels"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(signedToken string) (*sharedmodels.SignedDetails, error) {

	token, err := jwt.ParseWithClaims(signedToken, &sharedmodels.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sharedmodels.SECRET_KEY), nil
	})
	if err != nil && strings.Contains(err.Error(), "expired") {
		return nil, customErrors.ExpiredToken{}
	}
	if err != nil {
		return nil, customErrors.InvalidToken{}
	}

	claims, ok := token.Claims.(*sharedmodels.SignedDetails)
	if !ok {
		return nil, customErrors.InvalidToken{}
	}
	return claims, nil
}
