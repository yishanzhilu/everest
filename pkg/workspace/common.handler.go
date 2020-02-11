package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/http/server/middleware"
	"github.com/yishanzhilu/everest/pkg/models"
)

func abortWithPublicError(c *gin.Context, code int, err error, meta string) {
	if err != nil {
		cErr := c.Error(err).SetType(gin.ErrorTypePublic)
		if meta != "" {
			cErr.SetMeta(meta)
		}
		c.AbortWithStatusJSON(code, c.Errors.String())
		return
	}
}

func handleDBError(c *gin.Context, err error, meta string) {
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			abortWithPublicError(c, http.StatusNotFound, err, meta)
		} else {
			c.AbortWithError(http.StatusBadRequest, err).SetMeta(meta)
		}
		return
	}
}

// RegisterRoutes .
func RegisterRoutes(r *gin.RouterGroup) {
	r.Use(middleware.Authenticate)
	r.GET("/overview", getDoingGoalAndMissions)
	registerGoalRoutes(r)
	registerMissionRoutes(r)
	registerRecordRoutes(r)
	registerTodoRoutes(r)
	// r.GET("/missions", getMissionList)
	// r.GET("/mission/:id", getMission)
	// r.PATCH("/mission/:id", patchMission)
	// r.DELETE("/mission/:id", deleteMission)
}

func getDoingGoalAndMissions(c *gin.Context) {
	uid := c.MustGet(common.ContextUserID)
	var goals []models.GoalModel
	var missions []models.MissionModel
	err := common.MySQLClient.
		Preload("Missions", "status = ?", models.StatusDoing).
		Where("user_id = ? AND status = ?", uid, models.StatusDoing).
		Find(&goals).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	err = common.MySQLClient.Where(
		"user_id = ? AND status = ? AND goal_id = ?",
		uid,
		models.StatusDoing,
		0,
	).Find(&missions).
		Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{"goals": goals, "missions": missions})
}
