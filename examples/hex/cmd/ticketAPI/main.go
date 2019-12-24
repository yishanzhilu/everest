package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/yishanzhilu/api-template/examples/hex/internal/storage/mysql"
	redisRepo "github.com/yishanzhilu/api-template/examples/hex/internal/storage/redis"
	"github.com/yishanzhilu/api-template/examples/hex/internal/ticket"
)

func main() {
	db := mysqlConnect("user:pass@tcp(localhost:3306)/db_name?parseTime=true")
	defer db.Close()
	redis := redisConnect("localhost:6379", "foobared")
	defer redis.Close()

	ticketMysqlRepo := mysql.NewMysqlTicketRepository(db)
	ticketRepo := redisRepo.NewRedisModelRepository(redis, ticketMysqlRepo)
	ticketService := ticket.NewTicketService(ticketRepo)
	ticketHandler := ticket.NewTicketHandler(ticketService)

	router := gin.Default()
	router.GET("/tickets", ticketHandler.Get)
	router.GET("/ticket/{id}", ticketHandler.GetByID)
	router.POST("/tickets", ticketHandler.Create)
}

func redisConnect(url string, password string) *redis.Client {
	logrus.WithField("connection", url).Info("Connecting to Redis DB")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	})
	err := client.Ping().Err()

	if err != nil {
		logrus.Fatal(err)
	}
	return client
}

func mysqlConnect(url string) *gorm.DB {
	logrus.WithField("connection", url).Info("Connecting to MySQL DB")
	db, err := gorm.Open("mysql", url)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
