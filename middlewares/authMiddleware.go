package middlewares

import (
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/tudoresan/shared-finance-app-golang/customErrors"
	"github.com/TudorEsan/shared-finance-app-golang/releases/tag/v1.0.0/"
)

type SignedDetails struct {
	Email    string
	Username string
	Id       string
	jwt.StandardClaims
}

func GetSecretKey() string {
	return os.Getenv("SECRET_JWT")
}

var SECRET_KEY string = GetSecretKey()

func validateToken(signedToken string) (*SignedDetails, error) {

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil && strings.Contains(err.Error(), "expired") {
		return nil, customErrors.ExpiredToken{}
	}
	if err != nil {
		return nil, customErrors.InvalidToken{}
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, customErrors.InvalidToken{}
	}
	return claims, nil
}

func removeCookies(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", "", 60*60*24*30, "", "", false, false)
}

func VerifyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if token exists
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token Not Found"})
			removeCookies(c)
			c.Abort()
			return
		}
		// Check if Refresh Token exists
		_, err = c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh Token Not Found"})
			removeCookies(c)
			c.Abort()
			return
		}

		// Validate Token
		claims, err := validateToken(token)
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
