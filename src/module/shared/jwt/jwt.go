package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenMalformed occurs when the given access token is not valid JWT token
var ErrAccessTokenMalformed = fmt.Errorf("inavlid access token")

// AccessTokenExpired occurs when the given access token is expired
var ErrAccessTokenExpired = fmt.Errorf("inactive or expired token")

type CustomClaims struct {
	UserUID string `json:"user_uid"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userUID string) string {
	var ACCESS_TOKEN_EXPIRE_TIME, _ = strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRE_TIME"), 10, 64)
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

	now := time.Now()
	ttl := time.Duration(ACCESS_TOKEN_EXPIRE_TIME) * time.Second

	claims := CustomClaims{
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

// DecodeAccessToken decode the given token string into the token claims.
func DecodeAccessToken(accessToken string) (CustomClaims, error) {
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ACCESS_TOKEN_SECRET), nil
	})

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return CustomClaims{}, ErrAccessTokenMalformed
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return CustomClaims{}, ErrAccessTokenExpired
			} else {
				return CustomClaims{}, fmt.Errorf("error when decoding token")
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		// manually verify issuer
		// if !claims.VerifyIssuer(issuer, true) {
		// 	return nil, fmt.Errorf("invalid token issuer")
		// }

		return *claims, nil
	}
	return CustomClaims{}, fmt.Errorf("error when parsing into claims")
}

type TenantAccessTokenParam struct {
	UserUID string `json:"user_uid"`
	Domain  string `json:"domain"`
}

type TenantCustomClaims struct {
	UserUID string `json:"user_uid"`
	Domain  string `json:"domain"`
	jwt.RegisteredClaims
}

func GenerateTenantAccessToken(param TenantAccessTokenParam) string {
	var ACCESS_TOKEN_EXPIRE_TIME, _ = strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRE_TIME"), 10, 64)
	var ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")

	now := time.Now()
	ttl := time.Duration(ACCESS_TOKEN_EXPIRE_TIME) * time.Second

	claims := TenantCustomClaims{
		UserUID: param.UserUID,
		Domain:  param.Domain,
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
