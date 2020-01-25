package workspace

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func registerTaskRoutes(r *gin.RouterGroup) {
	r.POST("/tasks", postTask)
	r.GET("/tasks", getTaskList)
	r.GET("/task/:id", getTask)
	r.PATCH("/task/:id", patchTask)
	r.DELETE("/task/:id", deleteTask)
}

type postTaskBody struct {
	Content   string `json:"content" binding:"required,max=255"`
	Review    string `json:"review" binding:"max=255"`
	Minutes   uint16 `json:"minutes" binding:"max=480"`
	GoalID    uint64 `json:"goalID"`
	MissionID uint64 `json:"missionID"`
	Status    string `json:"status" binding:"oneof=todo done"`
}

func postTask(c *gin.Context) {
	var body postTaskBody
	var err error
	if err = c.BindJSON(&body); err != nil {
		return
	}
	if body.Status == "todo" && (body.Review != "" || body.Minutes > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "task with status todo can't have review or minutes"})
		return
	}

	uid := c.MustGet(common.ContextUserID).(uint64)
	task := models.TaskModel{
		Content: body.Content,
		Review:  body.Review,
		Minutes: body.Minutes,
		UserID:  uid,
		Status:  models.WorkStatsMap[body.Status],
	}

	if body.MissionID > 0 {
		mission, err := getMissionHelper(uid, body.MissionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"missionID": err.Error()})
			return
		}
		task.MissionID = mission.ID
		task.Mission.Title = mission.Title
		task.GoalID = mission.GoalID
		task.Goal.Title = mission.Goal.Title
	} else if body.GoalID > 0 {
		goal, err := getGoalHelper(uid, body.GoalID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goalID": err.Error()})
			return
		}
		task.GoalID = goal.ID
		task.Goal.Title = goal.Title
	}

	if err := common.MySQLClient.Create(&task).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// get user info for response
	// if err := common.MySQLClient.First(&task.User, task.UserID).Error; err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }

	c.JSON(http.StatusCreated, task)
}

func getTaskList(c *gin.Context) {
	goalIDStr := c.DefaultQuery("goalID", "0")
	goalID, err := strconv.ParseUint(goalIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	missionIDStr := c.DefaultQuery("missionID", "0")
	missionID, err := strconv.ParseUint(missionIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	status := c.DefaultQuery("status", "any")
	statusCode, ok := models.WorkStatsMap[status]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var tasks []models.TaskModel
	err = common.MySQLClient.
		Preload("Goal").
		Preload("Mission").
		Where(&models.TaskModel{
			UserID:    c.MustGet(common.ContextUserID).(uint64),
			Status:    statusCode,
			GoalID:    goalID,
			MissionID: missionID,
		}).
		Find(&tasks).Error
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func getTaskHelper(uid, taskID interface{}) (task *models.TaskModel, err error) {
	task = &models.TaskModel{}
	err = common.MySQLClient.
		Where("user_id = ?", uid).
		// Preload("User").
		Preload("Goal").
		Preload("Mission").
		First(task, taskID).Error
	return
}

func getTask(c *gin.Context) {
	mission, err := getTaskHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "task")
		return
	}
	c.JSON(http.StatusOK, *mission)
}

type patchTaskBody struct {
	Content string `json:"content" binding:"max=80"`
	Review  string `json:"review" binding:"max=255"`
	Minutes uint16 `json:"minutes" binding:"max=480"`
	Status  string `json:"status" binding:"omitempty,oneof=done"`
}

func patchTask(c *gin.Context) {

	var body patchTaskBody
	err := c.BindJSON(&body)
	if err != nil {
		return
	}
	// if status is not done,
	// user is update todo content
	// so, review and minutes can't be non-empty
	if body.Status == "" && (body.Review != "" || body.Minutes > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "can only add review when new task status is done"})
		return
	}

	newTask := models.TaskModel{
		Content: body.Content,
		Review:  body.Review,
		Status:  models.WorkStatsMap[body.Status],
		Minutes: body.Minutes,
	}

	task, err := getTaskHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "task")
		return
	}
	// finished task can't be updated
	if task.Status == models.StatusDone {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "finished task can't be updated"})
		return
	}
	err = common.MySQLClient.
		Model(task).
		Updates(newTask).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, *task)
}

func deleteTask(c *gin.Context) {
	task, err := getTaskHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "task")
		return
	}
	if err := common.MySQLClient.Delete(task).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Status(http.StatusOK)
}
