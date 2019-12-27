package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yishanzhilu/everest/pkg/common"
)

// WorkprofileHandler has workprofile GIN handlers, it is the adapter layer
type WorkprofileHandler interface {
	common.BaseHandler
	GetList(c *gin.Context)
	Create(c *gin.Context)
}

type workprofileHandler struct {
	workprofileService WorkprofileService
}

// NewWorkprofileHandler ...
func NewWorkprofileHandler(workprofileService WorkprofileService) WorkprofileHandler {
	return &workprofileHandler{
		workprofileService,
	}
}

func (h *workprofileHandler) RegisterPublicRoutes(r *gin.RouterGroup) {
	r.GET("/profile/:id", h.GetByID)
	r.GET("/profiles", h.GetList)
	r.POST("/profile", h.Create)
}

func (h *workprofileHandler) RegisterPrivateRoutes(r *gin.RouterGroup) {
}

func (h *workprofileHandler) GetList(c *gin.Context) {
	workprofile, _ := h.workprofileService.FindList()
	c.JSON(http.StatusOK, workprofile)
}

func (h *workprofileHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	workprofile, err := h.workprofileService.FindByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, workprofile)
}

func (h *workprofileHandler) Create(c *gin.Context) {
	var workprofile WorkprofileModel
	if err := c.BindJSON(&workprofile); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	if err := h.workprofileService.CreateWorkprofile(&workprofile); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	c.JSON(http.StatusOK, workprofile)
}
