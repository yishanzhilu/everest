package bootstrap

import (
	"time"

	"github.com/yishanzhilu/everest/pkg/common"
)

func initHTTPClient() {
	common.HTTPClient.SetTimeout(30 * time.Second)
}
