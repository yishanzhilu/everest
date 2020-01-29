package crypto

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTGuard will check visitor's identify and report to business handlers
type JWTGuard interface {
	CheckToken(tokenString string) (userID uint64, err error)
	SignToken(userID uint64, userName string) (string, error)
}

// jwtGuard implements JWTGuard
type jwtGuard struct {
	// Secret is JWT secret.
	secret   string
	duration int64
}

// NewJWTGuard JWTGuard.
func NewJWTGuard(secret string, duration time.Duration) JWTGuard {
	return &jwtGuard{secret, int64(duration)}
}

type everestClaims struct {
	UserID uint64 `json:"userID"`
	jwt.StandardClaims
}

// CheckToken will let JWTGuard check if a token is valid, and return visitor's userID
func (ja *jwtGuard) CheckToken(tokenString string) (userID uint64, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &everestClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ja.secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = errors.New("not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				err = errors.New("token expired")

			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				// Token is either expired or not active yet
				err = errors.New("token not valid yet")
			} else {
				err = fmt.Errorf("Couldn't handle this token: %e", err)
			}
		} else {
			err = fmt.Errorf("Couldn't handle this token: %e", err)
		}
		return
	}

	if claims, ok := token.Claims.(*everestClaims); ok {
		return claims.UserID, err
	}
	return 0, errors.New("Couldn't handle this token")
}

// SignToken implementation.
func (ja *jwtGuard) SignToken(userID uint64, userName string) (string, error) {
	// duration := ja.duration
	// if duration == 0 {
	// 	duration = 3600 // default 1 hour
	// }
	claims := everestClaims{
		userID,
		jwt.StandardClaims{
			// ExpiresAt: time.Now().Unix() + duration,
			Issuer: "everest",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(ja.secret))
	return ss, err
}

// stubJWTGuard implements JWTGuard as test stub
type stubJWTGuard struct {
	// Secret is JWT secret.
	secret   string
	duration int64
}

// NewStubJWtGuard implements JWTGuard as test stub
func NewStubJWtGuard() JWTGuard {
	return &stubJWTGuard{"", 1}
}

// CheckToken will let JWTGuard check if a token is valid, and return visitor's userID
func (s *stubJWTGuard) CheckToken(tokenString string) (userID uint64, err error) {
	if tokenString == "good token for stub guard" {
		return 1, nil
	}
	return 0, nil
}

// SignToken implementation.
func (s *stubJWTGuard) SignToken(userID uint64, userName string) (string, error) {
	return "", nil
}
