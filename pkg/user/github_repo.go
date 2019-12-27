package user

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yishanzhilu/everest/pkg/common"
	"gopkg.in/resty.v1"
)

// OauthRepo abstract github oauth repo layer logic
type OauthRepo interface {
	GetUserOauthToken(code string) (*GithubToken, error)
	GetUserOauthInfo(token string) (*GithubUser, error)
}

type githubRepo struct {
	restClient *resty.Client
	logger     *logrus.Logger
}

// NewGithubRepo create a new githubRepo instance
func NewGithubRepo(restClient *resty.Client, logger *logrus.Logger) OauthRepo {
	return &githubRepo{
		restClient,
		logger,
	}
}

func (r *githubRepo) GetUserOauthToken(code string) (*GithubToken, error) {
	r.logger.WithFields(logrus.Fields{
		"code": code,
	}).Debug("call start")

	var token = &GithubToken{}
	req := r.restClient.R().
		SetQueryParams(map[string]string{
			"code": code,
		}).
		SetResult(token)
	resp, err := req.
		Get("/login/oauth/access_token")
	r.logger.WithFields(logrus.Fields{
		"error": err,
		"resp":  resp,
	}).Debug("call done")

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() > 200 {
		return nil, errors.New(resp.String())
	}

	return token, nil
}

func (r githubRepo) GetUserOauthInfo(token string) (*GithubUser, error) {
	common.Logger.WithFields(logrus.Fields{
		"token": token,
	}).Debug("call start")
	gu := &GithubUser{}
	req := common.HTTPClient.R().
		SetHeaders(map[string]string{"Authorization": fmt.Sprintf("token %s", token)}).
		SetResult(gu)
	resp, err := req.
		Get("https://api.github.com/user")
	common.Logger.WithFields(logrus.Fields{
		"error": err,
		"resp":  resp,
	}).Debug("call done")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() > 200 {
		return nil, errors.New(resp.String())
	}
	return gu, nil
}
