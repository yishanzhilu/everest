package main

import (
	"github.com/yishanzhilu/api-template/pkg/bootstrap"
)

func main() {
	bootstrap.Boot()
	defer bootstrap.Close()
}
