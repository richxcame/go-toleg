package helpers

import (
	"gotoleg/web/entities"
	"gotoleg/web/middlewares"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	User entities.User `json:"user"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateJWT creates access and refresh tokens with user's username
func GenerateJWT(username string) (token Tokens, err error) {
	// Create access token
	accessTokenExp := time.Now().Add(3 * time.Hour)
	accessClaims := &Claims{
		User: entities.User{
			Username: username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: accessTokenExp},
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	token.AccessToken, err = accessToken.SignedString(JWT_SECRET)
	if err != nil {
		return Tokens{}, err
	}

	// Create refresh token
	refreshTokenExp := time.Now().Add(24 * time.Hour * 30)
	refreshClaims := &Claims{
		User: entities.User{
			Username: username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: refreshTokenExp},
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	token.RefreshToken, err = refreshToken.SignedString(JWT_SECRET)
	if err != nil {
		return Tokens{}, err
	}

	return token, nil
}

func RefreshToken(claims *middlewares.Claims) (token middlewares.Tokens, err error) {
	accessTokenTimeOut, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIMEOUT"))
	if err != nil {
		return middlewares.Tokens{}, err
	}
	expirationTime := time.Now().Add(time.Duration(accessTokenTimeOut) * time.Second)

	claims.ExpiresAt = &jwt.NumericDate{Time: expirationTime}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.AccessToken, err = accessToken.SignedString(middlewares.JWT_SECRET)
	if err != nil {
		return middlewares.Tokens{}, err
	}

	refreshTokenTimeOut, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIMEOUT"))
	if err != nil {
		return middlewares.Tokens{}, err
	}
	expirationTime = time.Now().Add(time.Duration(refreshTokenTimeOut) * time.Second)

	claims.ExpiresAt = &jwt.NumericDate{Time: expirationTime}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.RefreshToken, err = refreshToken.SignedString(middlewares.JWT_SECRET)
	if err != nil {
		return middlewares.Tokens{}, err
	}

	return token, nil
}
