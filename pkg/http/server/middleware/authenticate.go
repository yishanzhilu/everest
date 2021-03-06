package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yishanzhilu/everest/pkg/crypto"

	"github.com/sirupsen/logrus"
	"github.com/yishanzhilu/everest/pkg/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func extractToken(c *gin.Context) string {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		var err error
		tokenString, err = c.Cookie("x-tai-everest-token")
		if err != nil {
			return ""
		}
	}
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

// AssignGuard will check if token is valid,
// 		if true, then add KV authorized=true to gin.context. It will also add parsed JWT info to context
// 		if false, then add KV authorized=false to gin.context.
// Note: JWT won't terminate call handler if there is no token, it only do so when token is invalid
func AssignGuard(guard crypto.JWTGuard) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("authorized", false)
		tokenString := extractToken(c)
		if tokenString != "" {
			userID, err := guard.CheckToken(tokenString)
			if err != nil {
				logrus.Debug("Guard find BAD token")

				c.Set("tokenErr", err.Error())
			} else {
				logrus.Debug("Guard find GOOD token")
				c.Set("authorized", true)
			}
			c.Set("userID", userID)
			common.Logger.WithField("UserID", userID).Debug("AssignGuard report")
		}
		c.Next()
	}
}

// Authenticate will check if call is authorized based on gin.context info
func Authenticate(c *gin.Context) {
	authorized := c.GetBool("authorized")
	if !authorized {
		tokenErr := c.GetString("tokenErr")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized", "reason": tokenErr})
		return
	}
	c.Next()
}
