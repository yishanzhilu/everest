package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/heptiolabs/healthcheck"
	"github.com/yishanzhilu/api-template/pkg/common"
)

// RedisCheck 用于让 healthcheck 检查
func RedisCheck(cache *redis.Client) healthcheck.Check {
	return func() error {
		if _, err := cache.Ping().Result(); err != nil {
			return err
		}
		return nil
	}
}

func healthCheck(r *gin.Engine) {
	health := healthcheck.NewHandler()

	// db 链接
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(common.MySQLClient.DB(), 1*time.Second))
	// redis 链接
	health.AddReadinessCheck("redis", RedisCheck(common.RedisClient))

	// goroutine 限制
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	// 添加路由
	r.GET("/live", gin.WrapF(health.LiveEndpoint))
	r.GET("/ready", gin.WrapF(health.ReadyEndpoint))
}
