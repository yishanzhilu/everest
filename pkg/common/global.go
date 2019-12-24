package common

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// Logger is the common log agency in project
var Logger = logrus.New()

// MySQLClient is the common MySQL client in project
var MySQLClient *gorm.DB

// RedisClient is the common Redis client in project
var RedisClient *redis.Client

// HTTPClient is the common HTTP client in project
var HTTPClient *resty.Client = resty.New()

// JWTSecret 是 jwt token 的秘钥
var JWTSecret string
