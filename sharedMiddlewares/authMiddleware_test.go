package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	middlewares "github.com/TudorEsan/shared-finance-app-golang/sharedMiddlewares"
	sharedmodels "github.com/TudorEsan/shared-finance-app-golang/sharedModels"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func generateJwtToken(emailValidated bool, expirationTime time.Time) (string, error) {
	claims := &sharedmodels.SignedDetails{
		EmailValidated: emailValidated,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sharedmodels.SECRET_KEY)
	return tokenString, err
}

func TestVerifyAuth(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.VerifyAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerifyAuthWithToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.VerifyAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "token"})

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerifyAuthWithTokenAndRefreshToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.VerifyAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "token"})
	req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "refreshToken"})

	// Test
	router.ServeHTTP(w, req)

	t.Log(w.Body.String())
	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestEmailNotVerified(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.VerifyAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	token, _ := generateJwtToken(false, time.Now().Add(time.Hour))

	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	req.AddCookie(&http.Cookie{Name: "refreshToken", Value: token})

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "{\"message\":\"Email Not Validated\"}", w.Body.String())
}

func TestExpiredToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.VerifyAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	token, _ := generateJwtToken(true, time.Now().Add(-time.Hour))
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	req.AddCookie(&http.Cookie{Name: "refreshToken", Value: token})

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "{\"message\":\"Token Expired\"}", w.Body.String())
}

func TestValidToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.VerifyAuth())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	token, _ := generateJwtToken(true, time.Now().Add(time.Hour))
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	req.AddCookie(&http.Cookie{Name: "refreshToken", Value: token})

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}
