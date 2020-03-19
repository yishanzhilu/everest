package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func registerGoalRoutes(r *gin.RouterGroup) {
	r.POST("/goals", postGoal)
	r.GET("/goals", getGoalList)
	r.GET("/goal/:id", getGoal)
	r.PATCH("/goal/:id", patchGoal)
	r.DELETE("/goal/:id", deleteGoal)
}

type postGoalBody struct {
	Title       string `json:"title" binding:"required,max=80"`
	Description string `json:"description" binding:"max=1000"`
	Status      string `json:"status" binding:"oneof=doing todo done drop"`
}

func postGoal(c *gin.Context) {
	var body postGoalBody
	if err := c.BindJSON(&body); err != nil {
		return
	}
	goal := models.GoalModel{
		Title:       body.Title,
		Description: body.Description,
		UserID:      c.MustGet(common.ContextUserID).(uint64),
		Status:      models.WorkStatsMap[body.Status],
	}
	if err := common.MySQLClient.Create(&goal).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// get user info for response
	// if err := common.MySQLClient.First(&goal.User, goal.UserID).Error; err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }

	c.JSON(http.StatusCreated, goal)
}

func getGoalList(c *gin.Context) {
	status := c.DefaultQuery("status", "any")
	withMission := c.DefaultQuery("missions", "false")
	statusCode, ok := models.WorkStatsMap[status]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var goals []models.GoalModel
	db := common.MySQLClient.
		Where("user_id = ? ", c.MustGet(common.ContextUserID))
		// Preload("Missions", "status = ?", models.StatusDoing)
	if statusCode != models.StatusAny {
		db = db.Where("status = ?", statusCode)
	}
	if withMission == "true" {
		db = db.Preload("Missions")

	}
	err := db.Order("updated_at desc").Find(&goals).Error
	if err != nil {
		handleDBError(c, err, "goal")
		return
	}
	c.JSON(http.StatusOK, goals)
}

func getGoalHelper(uid, goalID interface{}) (goal *models.GoalModel, err error) {
	goal = &models.GoalModel{}
	err = common.MySQLClient.
		Where("user_id = ?", uid.(uint64)).
		// Preload("User").
		// Preload("Missions", "status = ?", models.StatusDoing).
		First(goal, goalID).Error
	return
}

func getGoal(c *gin.Context) {
	goal, err := getGoalHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "goal")
		return
	}
	c.JSON(http.StatusOK, *goal)
}

type patchGoalBody struct {
	Title       string `json:"title" binding:"max=80"`
	Description string `json:"description" binding:"max=1000"`
	Status      string `json:"status" binding:"omitempty,oneof=doing todo done drop"`
}

func patchGoal(c *gin.Context) {

	var body patchGoalBody
	err := c.BindJSON(&body)
	if err != nil {
		return
	}
	newGoal := models.GoalModel{
		Title:       body.Title,
		Description: body.Description,
	}
	if body.Status != "" {
		newGoal.Status = models.WorkStatsMap[body.Status]
	}

	goal, err := getGoalHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "goal")
		return
	}
	err = common.MySQLClient.
		Model(goal).
		Updates(newGoal).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, *goal)
}

func deleteGoal(c *gin.Context) {
	goal, err := getGoalHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "goal")
		return
	}
	err = common.MySQLClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(goal).Error; err != nil {
			return err
		}

		if err := tx.Where("goal_id = ?", goal.ID).
			Delete(models.RecordModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("goal_id = ?", goal.ID).
			Delete(models.MissionModel{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
