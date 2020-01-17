package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/crypto"
	"github.com/yishanzhilu/everest/pkg/http/server/middleware"
	"github.com/yishanzhilu/everest/pkg/models"
)

// Handler has user GIN handlers, it is the adapter layer
type Handler interface {
	common.BaseHandler
	OauthGithub(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetAuthenticatedUser(c *gin.Context)
	UpdateAuthenticatedUser(c *gin.Context)
}

type handler struct {
	service Service
	guard   crypto.JWTGuard
}

// NewHandler ...
func NewHandler(service Service, guard crypto.JWTGuard) Handler {
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
	r.Use(middleware.Authenticate())
	r.GET("", h.GetAuthenticatedUser)
	r.PATCH("", h.UpdateAuthenticatedUser)
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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := h.guard.SignToken(user.ID, user.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "refreshToken": user.RefreshToken})
}

type refreshTokenParameters struct {
	UserID       uint64 `json:"userID"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshToken implementation.
func (h *handler) RefreshToken(c *gin.Context) {
	// get userID in context which was set in AssignGuard middleware
	rp := &refreshTokenParameters{}
	if err := c.ShouldBindJSON(rp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	common.Logger.Debug("EverestUserID\t", rp.UserID, "\tEverestRefreshToken\t", rp.RefreshToken)
	if rp.UserID == 0 || rp.RefreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid token or refresh token"})
		return
	}
	user, err := h.service.GetByID(rp.UserID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if user.RefreshToken != rp.RefreshToken {
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

type updateAuthenticatedUserParameters struct {
	Name      string `binding:"max=80"`
	AvatarURL string `json:"avatarUrl" binding:"url"`
}

// UpdateAuthenticatedUser implementation.
func (h *handler) UpdateAuthenticatedUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	up := &updateAuthenticatedUserParameters{}
	if err := c.ShouldBindJSON(up); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	newUser := &models.UserModel{}
	newUser.Name = up.Name
	newUser.AvatarURL = up.AvatarURL
	user, err := h.service.UpdateByID(userID, newUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetAuthenticatedUser implementation.
func (h *handler) GetAuthenticatedUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)
	user, err := h.service.GetByID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
