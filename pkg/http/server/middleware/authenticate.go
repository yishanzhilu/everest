package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yishanzhilu/api-template/pkg/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func extractToken(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
	return tokenString
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func validateToken(token *jwt.Token, c *gin.Context) error {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["sub"].(string)
		c.Set("userID", userID)
		logrus.WithField("UserID", userID).Info("Validated user")
		return nil
	}
	return errors.New("Invalid Token")
}

// Authenticate 检查请求是否登录
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid token"})
		}
		token, err := parseToken(tokenString)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid token"})
			return
		}
		err = validateToken(token, c)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid token"})
			return
		}
	}
}
