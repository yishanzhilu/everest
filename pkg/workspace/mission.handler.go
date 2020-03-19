package workspace

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func registerMissionRoutes(r *gin.RouterGroup) {
	r.POST("/missions", postMission)
	r.GET("/missions", getMissionList)
	r.GET("/mission/:id", getMission)
	r.PATCH("/mission/:id", patchMission)
	r.DELETE("/mission/:id", deleteMission)
}

type postMissionBody struct {
	Title       string `json:"title" binding:"required,max=80"`
	Description string `json:"description" binding:"max=1000"`
	GoalID      uint64 `json:"goalID"`
	Status      string `json:"status" binding:"oneof=doing todo done drop"`
}

func postMission(c *gin.Context) {
	var body postMissionBody
	var err error
	if err = c.BindJSON(&body); err != nil {
		return
	}
	uid := c.MustGet(common.ContextUserID).(uint64)
	mission := models.MissionModel{
		Title:       body.Title,
		Description: body.Description,
		UserID:      uid,
		Status:      models.WorkStatsMap[body.Status],
	}
	if body.GoalID > 0 {
		goal, err := getGoalHelper(uid, body.GoalID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goalID": err.Error()})
			return
		}
		mission.GoalID = goal.ID
		mission.Goal.ID = goal.ID
		mission.Goal.Title = goal.Title
	}

	if err := common.MySQLClient.
		Create(&mission).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// get user info for response
	if err := common.MySQLClient.First(&mission.User, mission.UserID).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, mission)
}

func getMissionList(c *gin.Context) {
	status := c.DefaultQuery("status", "any")
	goalIDStr := c.DefaultQuery("goalID", "0")
	goalID, err := strconv.ParseUint(goalIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	statusCode, ok := models.WorkStatsMap[status]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var missions []models.MissionModel
	err = common.MySQLClient.
		Where(&models.MissionModel{
			UserID: c.MustGet(common.ContextUserID).(uint64),
			Status: statusCode,
		}).
		Where("goal_id = ?", goalID).
		Preload("Goal").
		Find(&missions).Error
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, missions)
}

func getMissionHelper(uid, missionID interface{}) (mission *models.MissionModel, err error) {
	mission = &models.MissionModel{}
	err = common.MySQLClient.
		Where("user_id = ?", uid).
		// Preload("User").
		Preload("Goal").
		First(mission, missionID).Error
	return
}

func getMission(c *gin.Context) {
	mission, err := getMissionHelper(c.MustGet(common.ContextUserID).(uint64), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "mission")
		return
	}
	c.JSON(http.StatusOK, *mission)
}

type patchMissionBody struct {
	Title       string  `json:"title" binding:"max=80"`
	Description string  `json:"description" binding:"max=1000"`
	Status      string  `json:"status" binding:"omitempty,oneof=doing todo done drop"`
	GoalIDPtr   *uint64 `json:"goalID,omitempty"`
}

func patchMission(c *gin.Context) {

	var body patchMissionBody
	err := c.BindJSON(&body)
	if err != nil {
		return
	}
	// newMission := models.GoalModel{
	// 	Title:       body.Title,
	// 	Description: body.Description,
	// }
	// if body.Status != "" {
	// 	newMission.Status = models.WorkStatsMap[body.Status]
	// }
	uid := c.MustGet(common.ContextUserID)
	mission, err := getMissionHelper(uid, c.Param("id"))
	if err != nil {
		handleDBError(c, err, "mission")
		return
	}
	newMission := make(map[string]interface{})
	if body.Title != "" && body.Title != mission.Title {
		newMission["title"] = body.Title
	}
	if body.Description != "" && body.Description != mission.Description {
		newMission["description"] = body.Description
	}
	newStatus := models.WorkStatsMap[body.Status]
	if body.Status != "" && mission.Status != newStatus {
		newMission["status"] = newStatus
	}
	if body.GoalIDPtr != nil {
		goalID := *(body.GoalIDPtr)
		if mission.GoalID != goalID {
			if goalID > 0 {
				goal, err := getGoalHelper(uid, goalID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goalID": err.Error()})
					return
				}
				mission.Goal.ID = goal.ID
				mission.Goal.Title = goal.Title
				newMission["goal_id"] = goal.ID
			} else {
				mission.Goal.ID = 0
				mission.Goal.Title = ""
				newMission["goal_id"] = 0
			}
		}
	}

	err = common.MySQLClient.
		Model(mission).
		Updates(newMission).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, *mission)
}

func deleteMission(c *gin.Context) {
	mission, err := getMissionHelper(c.MustGet(common.ContextUserID).(uint64), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "mission")
		return
	}

	err = common.MySQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(mission).Error; err != nil {
			return err
		}

		if err := tx.Where("mission_id = ?", mission.ID).
			Delete(models.RecordModel{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		handleDBError(c, err, "mission")
		return
	}
	c.Status(http.StatusOK)
}
