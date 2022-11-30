package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"

	"final-project-backend/config"
	"final-project-backend/internal/dtos"
	"final-project-backend/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func validateToken(encodedToken string, secret string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(secret), nil
	})
}

func validateAccessToken(encodedToken string) (*jwt.Token, error) {
	jwtEnv := config.InitConfigJwt()

	return validateToken(encodedToken, jwtEnv[2])
}

func AuthorizeJWT(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		helpers.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}

	tokenString := authHeader[7:]

	token, err := validateAccessToken(tokenString)
	if err != nil {
		helpers.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}

	userJson, _ := json.Marshal(claims["data"])
	jwtData := dtos.JwtData{}
	err = json.Unmarshal(userJson, &jwtData)
	if err != nil {
		helpers.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}
	c.Set("user", jwtData)
	c.Next()

}
