package main

import (
	"github.com/yishanzhilu/everest/pkg/bootstrap"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func main() {
	bootstrap.Boot()
	defer bootstrap.Cleanup()
	models.Migrate(common.MySQLClient)
}
