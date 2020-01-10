package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/crypto"
)

// Handler has user GIN handlers, it is the adapter layer
type Handler interface {
	common.BaseHandler
	OauthGithub(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type handler struct {
	service Service
	guard   *crypto.JWTGuard
}

// NewHandler ...
func NewHandler(service Service, guard *crypto.JWTGuard) Handler {
	return &handler{
		service,
		guard,
	}
}

func (h *handler) RegisterPublicRoutes(r *gin.RouterGroup) {
	r.GET("/oauth/github", h.OauthGithub)
	r.POST("/token", h.RefreshToken)
}

func (h *handler) RegisterPrivateRoutes(r *gin.RouterGroup) {
}

// OauthGithub will use github oauth code to find user in yishan db,
// if user not exist, it will creat user, finally, it will return
// a JWT token with user info included
func (h *handler) OauthGithub(c *gin.Context) {
	var code string
	var ok bool
	if code, ok = c.GetQuery("code"); !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user, err := h.service.GetOrCreateUserWithGithubOauth(code)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Oauth fail"})
		return
	}
	token, err := h.guard.SignToken(user.ID, user.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "refreshToken": user.RefreshToken})
}

// RefreshToken implementation.
func (h *handler) RefreshToken(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)
	refreshToken := c.GetHeader("RefreshToken")
	common.Logger.Debug("userID\t", userID, "\tRefreshToken\t", refreshToken)
	if userID == 0 || refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid token or refresh token"})
		return
	}
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if user.RefreshToken != refreshToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid refreshToken"})
		return
	}
	token, err := h.guard.SignToken(user.ID, user.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
