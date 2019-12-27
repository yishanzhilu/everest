package user

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yishanzhilu/everest/pkg/crypto"
	"gopkg.in/resty.v1"
)

// ConstructNewUserHandler ..
func ConstructNewUserHandler(logger *logrus.Logger, db *gorm.DB, guard *crypto.JWTGuard) Handler {
	githubClient := resty.New().
		SetTimeout(30*time.Second).
		SetHostURL("https://github.com").
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"client_id":     viper.GetString("github.client.id"),
			"client_secret": viper.GetString("github.client.secret"),
		})
	gr := NewGithubRepo(githubClient, logger)
	mr := NewMysqlUserRepository(db)
	userService := NewUserService(mr, gr)
	userHandler := NewHandler(userService, guard)
	return userHandler
}
