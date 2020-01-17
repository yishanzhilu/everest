package user_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yishanzhilu/everest/pkg/crypto"
	"github.com/yishanzhilu/everest/pkg/http/server/middleware"
	. "github.com/yishanzhilu/everest/pkg/user"
)

var _ = Describe("User Handler", func() {
	gin.SetMode(gin.TestMode)

	var (
		h    Handler
		c    *gin.Context
		resp *httptest.ResponseRecorder
		r    *gin.Engine
	)
	BeforeEach(func() {
		guard := crypto.NewStubJWtGuard()
		h = NewHandler(NewStubService(), guard)
		resp = httptest.NewRecorder()
		c, r = gin.CreateTestContext(resp)
		r.Use(middleware.AssignGuard(guard))
		g := r.Group("user")

		h.RegisterPublicRoutes(g)
		h.RegisterPrivateRoutes(g)
	})
	Context("when update authenticated user", func() {
		It("should return 404 if use post", func() {
			c.Request, _ = http.NewRequest(http.MethodPost, "/user", nil)
			r.ServeHTTP(resp, c.Request)
			Expect(resp.Code).To(Equal(http.StatusNotFound))
		})
		It("should return 403 if has no auth key", func() {
			c.Request, _ = http.NewRequest(http.MethodPatch, "/user", nil)
			r.ServeHTTP(resp, c.Request)
			Expect(resp.Code).To(Equal(http.StatusForbidden))
		})
		It("should return 400 if has no patch body in json", func() {
			c.Set("authorized", true)
			c.Request, _ = http.NewRequest(http.MethodPatch, "/user", nil)
			c.Request.Header.Set("Authorization", "good token for stub guard")
			r.ServeHTTP(resp, c.Request)
			logrus.Infoln(resp.Body.String())
			Expect(resp.Code).To(Equal(422))
		})
	})
})
