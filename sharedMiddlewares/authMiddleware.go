package middlewares

import (
	"net/http"
	"strings"

	"github.com/TudorEsan/shared-finance-app-golang/customErrors"
	sharedmodels "github.com/TudorEsan/shared-finance-app-golang/sharedModels"
	"github.com/gin-gonic/gin"
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

func RemoveCookies(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", "", 60*60*24*30, "", "", false, false)
}

func VerifyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if token exists
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token Not Found"})
			RemoveCookies(c)
			c.Abort()
			return
		}
		// Check if Refresh Token exists
		_, err = c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh Token Not Found"})
			RemoveCookies(c)
			c.Abort()
			return
		}

		// Validate Token
		claims, err := ValidateToken(token)
		switch e := err.(type) {

		case nil:
			// token ok -> user authorized
			c.Set("UserId", claims.Id)
			c.Next()
			return
		case *customErrors.ExpiredToken:
			// Token expired -> client should refresh the tokens
			c.JSON(http.StatusInternalServerError, gin.H{"message": "token expired"})
			c.Abort()
			return
		default:
			// Token invalid or any other error-> reject Request
			c.JSON(http.StatusUnauthorized, customErrors.GetJsonError(e))
			c.Abort()
			return
		}
	}
}
