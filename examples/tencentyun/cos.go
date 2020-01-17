package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/yishanzhilu/everest/lib/tencentyun/sts"
	"github.com/yishanzhilu/everest/pkg/common"
)

func main() {
	getSts()
}

func upload() {
	u, _ := url.Parse("https://examplebucket-1250000000.cos.ap-shanghai.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  "SecretID",
			SecretKey: "SecretKey",
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	name := "test/objectPut.go"
	// 1. 通过字符串上传对象
	f := strings.NewReader("test")

	_, err := c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}
}

func getSts() {
	viper.SetConfigName("viper.local")           // name of config file (without extension)
	viper.AddConfigPath("./configs/")            // for dev, from bin/
	viper.AddConfigPath("../../configs/")        // for test, from pkg/bootstrap/
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		common.Logger.Error("配置初始化失败")
		panic(err)
	}
	appID := viper.GetString("tencentyun.app-id")
	bucket := viper.GetString("tencentyun.avatar-bucket.name")
	c := sts.NewClient(
		viper.GetString("tencentyun.secret-id"),
		viper.GetString("tencentyun.secret-key"),
		nil,
	)
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          viper.GetString("tencentyun.avatar-bucket.region"),
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:" + viper.GetString("tencentyun.avatar-bucket.region") + ":uid/" + appID + ":" + bucket + "/1",
					},
				},
			},
		},
	}
	res, err := c.GetCredential(opt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
	fmt.Printf("SessionToken	%+v\n\nTmpSecretID	%+v\n\nTmpSecretKey	%+v\n", res.Credentials.SessionToken, res.Credentials.TmpSecretID, res.Credentials.TmpSecretKey)
}
