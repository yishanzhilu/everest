package user_test

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/yishanzhilu/everest/pkg/user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/resty.v1"
)

var _ = Describe("GithubRepo", func() {
	var (
		server       *ghttp.Server
		githubClient *resty.Client
		githubRepo   user.OauthRepo
		logger       = logrus.New()
	)
	logger.SetLevel(logrus.DebugLevel)

	AfterEach(func() {
		server.Close()
	})
	BeforeEach(func() {
		server = ghttp.NewServer()
		githubClient = resty.New().SetHostURL(server.URL())
		githubRepo = user.NewGithubRepo(githubClient, logger)
	})
	Context("GetUserOauthToken", func() {

		It("should success with valid code", func() {
			server.RouteToHandler("GET", "/login/oauth/access_token",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/login/oauth/access_token", "code=test-code"),
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						user.NewGithubTokenResp(
							"test-token",
							"",
							"Bearer",
						),
					),
				),
			)
			token, err := githubRepo.GetUserOauthToken("test-code")
			Ω(token).ShouldNot(BeNil())
			Ω(token.AccessToken).Should(Equal("test-token"))
			Ω(err).Should(BeNil())

		})

		It("should success with bad code", func() {
			server.RouteToHandler("GET", "/login/oauth/access_token",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/login/oauth/access_token", "code="),
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						user.NewGithubTokenErr(),
					),
				),
			)
			token, err := githubRepo.GetUserOauthToken("")
			Ω(token).ShouldNot(BeNil())
			Ω(token.ErrorMessage).Should(Equal("bad_verification_code"))
			Ω(err).Should(BeNil())
		})
	})
})
