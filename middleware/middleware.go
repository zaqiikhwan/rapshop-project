package middleware

import (
	"errors"
	"net/http"
	"os"
	"rapsshop-project/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "auth header is not valid", errors.New("bearer header auth not found"))
			c.Abort()
			return
		}

		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := ValidateToken(tokenString)
		if err != nil {
			utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "token is not valid", err)
			c.Abort()
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "token or type is not valid", err)
			c.Abort()
			return
		}

		c.Set("id", claim["id"])
	}
}

func ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GenerateToken(id string) (string, error) {
	claim := jwt.MapClaims{}
	claim["id"] = id

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
