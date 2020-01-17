package tencentyun

import (
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// ObjectStorageRepo adapt to cloud storage services
type ObjectStorageRepo interface {
}

type tencentObjectStorageRepo struct {
	client *cos.Client
}

// NewObjectStorageRepo .
func NewObjectStorageRepo(rawurl, id, key string) ObjectStorageRepo {
	u, _ := url.Parse(rawurl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  id,
			SecretKey: key,
		},
	})
	return &tencentObjectStorageRepo{
		client,
	}
}
