package middlewares

import (
	"gotoleg/web/entities"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	User entities.User `json:"user"`
	jwt.RegisteredClaims
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := &Claims{}
		var token string

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "token_required",
				"message": "Auth token is required",
			})
			return
		}
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) > 1 {
			token = splitToken[1]
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "token_wrong",
				"message": "Invalid token",
			})
			return
		}
		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return JWT_SECRET, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Couldn't parse token",
			})
			return
		}

		if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "token_expired",
				"message": "Token expired",
			})
			return
		}

		if !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_token",
				"message": "Invalid token",
			})
			return
		}
		c.Next()
	}
}
