package helpers

import (
	"strconv"
	"time"

	"final-project-backend/config"

	"github.com/golang-jwt/jwt/v4"
)

type IdTokenClaims struct {
	jwt.RegisteredClaims
	Data interface{} `json:"data"`
}

func GenerateJwtToken(data interface{}, expiredMinute int64, issuer string, secret string) (*string, error) {
	idExp := expiredMinute * 60
	unixTime := time.Now().Unix()
	tokenExp := unixTime + idExp

	claims := &IdTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
			ExpiresAt: &jwt.NumericDate{
				Time: time.Unix(tokenExp, 0),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		Data: data,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func GenerateAccessToken(data interface{}) (*string, error) {
	jwtEnv := config.InitConfigJwt()
	expired, _ := strconv.ParseInt(jwtEnv[1], 10, 64)

	return GenerateJwtToken(data, expired, jwtEnv[0], jwtEnv[2])
}

func GenerateRefreshToken(data interface{}) (*string, error) {
	jwtEnv := config.InitConfigJwt()
	expired, _ := strconv.ParseInt(jwtEnv[3], 10, 64)

	return GenerateJwtToken(data, expired, jwtEnv[0], jwtEnv[4])
}
