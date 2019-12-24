package main

import (
	"github.com/yishanzhilu/everest/pkg/bootstrap"
)

func main() {
	bootstrap.Boot()
	defer bootstrap.Close()
}
