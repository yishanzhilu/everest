package bootstrap

import (
	"time"

	"github.com/yishanzhilu/api-template/pkg/common"
)

func initHTTPClient() {
	common.HTTPClient.SetTimeout(30 * time.Second)
}
