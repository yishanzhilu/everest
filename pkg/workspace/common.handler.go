package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/http/server/middleware"
)

func abortWithPublicError(c *gin.Context, code int, err error, meta string) {
	if err != nil {
		cErr := c.Error(err).SetType(gin.ErrorTypePublic)
		if meta != "" {
			cErr.SetMeta(meta)
		}
		c.AbortWithStatusJSON(code, c.Errors.JSON())
		return
	}
}

func handleDBError(c *gin.Context, err error, meta string) {
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			abortWithPublicError(c, http.StatusNotFound, err, meta)
		} else {
			c.AbortWithError(http.StatusBadRequest, c.Error(err).SetMeta(meta))
		}
		return
	}
}

// RegisterRoutes .
func RegisterRoutes(r *gin.RouterGroup) {
	r.Use(middleware.Authenticate)
	registerGoalRoutes(r)
	registerMissionRoutes(r)
	registerTaskRoutes(r)
	// r.GET("/missions", getMissionList)
	// r.GET("/mission/:id", getMission)
	// r.PATCH("/mission/:id", patchMission)
	// r.DELETE("/mission/:id", deleteMission)
}
