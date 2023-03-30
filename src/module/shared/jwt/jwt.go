package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type customClaims struct {
	UserUID string `json:"user_uid"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userUID string) string {
	var ACCESS_TOKEN_EXPIRE_TIME, _ = strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRE_TIME"), 10, 64)
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

	now := time.Now()
	ttl := time.Duration(ACCESS_TOKEN_EXPIRE_TIME) * time.Second

	claims := customClaims{
		UserUID: userUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(ttl),
			},
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(ACCESS_TOKEN_SECRET))
	return ss
}

func DecodeAccessToken(accessToken string) (customClaims, error) {
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	token, err := jwt.ParseWithClaims(accessToken, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ACCESS_TOKEN_SECRET), nil
	})

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return customClaims{}, fmt.Errorf("token tidak valid")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return customClaims{}, fmt.Errorf("inactive or expired token")
			} else {
				return customClaims{}, fmt.Errorf("error when decoding token")
			}
		}
	}

	if claims, ok := token.Claims.(*customClaims); ok {
		// manually verify issuer
		// if !claims.VerifyIssuer(issuer, true) {
		// 	return nil, fmt.Errorf("invalid token issuer")
		// }

		return *claims, nil
	}
	return customClaims{}, fmt.Errorf("error when parsing into claims")
}
