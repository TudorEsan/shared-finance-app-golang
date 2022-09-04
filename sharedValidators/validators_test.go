package sharedvalidators

import (
	"testing"
	"time"

	sharedmodels "github.com/TudorEsan/shared-finance-app-golang/sharedModels"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt"
)

func generateJwtToken(valid bool) (string, error) {
	var expirationTime time.Time
	if valid {
		expirationTime = time.Now().Add(time.Second * 10)
	} else {
		expirationTime = time.Now().Add(time.Second * -10)
	}
	claims := &sharedmodels.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sharedmodels.SECRET_KEY)
	return tokenString, err
}

func TestValidateToken(t *testing.T) {

	t.Run("Valid token test", func(t *testing.T) {
		token, err := generateJwtToken(true)
		if err != nil {
			t.Error(err)
		}
		_, err = ValidateToken(token)
		assert.Equal(t, err, nil)
	})

	t.Run("Unvalid token test", func(t *testing.T) {
		token, err := generateJwtToken(false)
		if err != nil {
			t.Error(err)
		}
		_, err = ValidateToken(token)
		assert.NotEqual(t, err, nil)
	})

}
