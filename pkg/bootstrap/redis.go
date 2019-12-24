package bootstrap

import (
	"github.com/go-redis/redis"
	"github.com/yishanzhilu/api-template/pkg/common"
)

func mustConnectRedis(url string, password string) {
	common.Logger.WithField("connection", url).Info("Connecting to Redis DB")
	client := redis.NewClient(&redis.Options{
		Addr:       url,
		Password:   password,
		DB:         0,
		MaxRetries: 2,
	})
	err := client.Ping().Err()
	if err != nil {
		common.Logger.WithField("error", err).Fatal("Failed to connect Redis DB")
	}
	common.Logger.Info("Connect to Redis DB success")
	common.RedisClient = client
}
