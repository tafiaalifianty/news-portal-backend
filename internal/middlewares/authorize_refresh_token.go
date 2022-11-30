package middlewares

import (
	"encoding/json"
	"net/http"

	"final-project-backend/config"
	"final-project-backend/internal/dtos"
	"final-project-backend/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func validateRefreshToken(encodedToken string) (*jwt.Token, error) {
	jwtEnv := config.InitConfigJwt()

	return validateToken(encodedToken, jwtEnv[4])
}

func AuthorizeRefreshToken(c *gin.Context) {
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

	token, err := validateRefreshToken(tokenString)
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
	refreshTokenData := dtos.RefreshTokenData{}
	err = json.Unmarshal(userJson, &refreshTokenData)
	if err != nil {
		helpers.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}
	c.Set("refresh_token", dtos.RefreshTokenContext{
		ID:           refreshTokenData.ID,
		RefreshToken: tokenString,
	})
	c.Next()

}
