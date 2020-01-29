package workspace

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yishanzhilu/everest/pkg/common"
	"github.com/yishanzhilu/everest/pkg/models"
)

func registerTodoRoutes(r *gin.RouterGroup) {
	r.POST("/todos", postTodo)
	r.GET("/todos", getTodoList)
	r.GET("/todo/:id", getTodo)
	r.PATCH("/todo/:id", patchTodo)
	r.DELETE("/todo/:id", deleteTodo)
}

type postTodoBody struct {
	Content   string `json:"content" binding:"required,max=255"`
	Minutes   uint16 `json:"minutes" binding:"max=480"`
	GoalID    uint64 `json:"goalID"`
	MissionID uint64 `json:"missionID"`
}

func postTodo(c *gin.Context) {
	var body postTodoBody
	var err error
	if err = c.BindJSON(&body); err != nil {
		return
	}

	uid := c.MustGet(common.ContextUserID).(uint64)
	todo := models.TodoModel{
		Content: body.Content,
		UserID:  uid,
	}
	todo.Status = models.StatusTodo

	if body.MissionID > 0 {
		mission, err := getMissionHelper(uid, body.MissionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"missionID": err.Error()})
			return
		}
		todo.MissionID = mission.ID
		todo.Mission.ID = mission.ID
		todo.Mission.Title = mission.Title
		todo.GoalID = mission.GoalID
		todo.Goal.ID = mission.Goal.ID
		todo.Goal.Title = mission.Goal.Title
	} else if body.GoalID > 0 {
		goal, err := getGoalHelper(uid, body.GoalID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goalID": err.Error()})
			return
		}
		todo.GoalID = goal.ID
		todo.Goal.ID = goal.ID
		todo.Goal.Title = goal.Title
	}

	if err := common.MySQLClient.Create(&todo).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func getTodoList(c *gin.Context) {
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

	status := c.DefaultQuery("status", "todo")
	statusCode, ok := models.WorkStatsMap[status]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var todos []models.TodoModel

	query := &models.TodoModel{
		UserID:    c.MustGet(common.ContextUserID).(uint64),
		GoalID:    goalID,
		MissionID: missionID,
	}
	query.Status = statusCode
	err = common.MySQLClient.
		Preload("Goal").
		Order("id").
		Preload("Mission").
		Where(query).
		Find(&todos).Error
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

func getTodoHelper(uid, todoID interface{}) (todo *models.TodoModel, err error) {
	todo = &models.TodoModel{}
	err = common.MySQLClient.
		Where("user_id = ?", uid).
		Preload("Goal").
		Preload("Mission").
		First(todo, todoID).Error
	return
}

func getTodo(c *gin.Context) {
	todo, err := getTodoHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "todo")
		return
	}
	c.JSON(http.StatusOK, *todo)
}

type patchTodoBody struct {
	Content string `json:"content" binding:"max=80"`
	Status  string `json:"status" binding:"omitempty,oneof=done"`
	// GoalID is a pointer so we can know if user pass 0 explictly which means relove relation
	GoalIDPtr *uint64 `json:"goalID,omitempty"`
	MissionID uint64  `json:"missionID"`
}

func patchTodo(c *gin.Context) {

	var body patchTodoBody
	err := c.BindJSON(&body)
	if err != nil {
		return
	}
	uid := c.MustGet(common.ContextUserID)

	todo, err := getTodoHelper(uid, c.Param("id"))
	if err != nil {
		handleDBError(c, err, "task")
		return
	}

	newTodo := make(map[string]interface{})
	if body.Content != "" && body.Content != todo.Content {
		newTodo["content"] = body.Content
	}
	newStatus := models.WorkStatsMap[body.Status]
	if body.Status != "" && todo.Status != newStatus {
		newTodo["status"] = newStatus
	}
	if body.MissionID > 0 && todo.MissionID != body.MissionID {
		mission, err := getMissionHelper(uid, body.MissionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"missionID": err.Error()})
			return
		}
		todo.MissionID = mission.ID
		todo.Mission.Title = mission.Title
		newTodo["mission_id"] = mission.ID
		todo.GoalID = mission.Goal.ID
		todo.Goal.Title = mission.Goal.Title
		newTodo["goal_id"] = mission.GoalID

	} else if body.GoalIDPtr != nil {
		goalID := *(body.GoalIDPtr)
		if todo.GoalID != goalID {
			if goalID > 0 {
				goal, err := getGoalHelper(uid, goalID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"goalID": err.Error()})
					return
				}
				todo.Goal.ID = goal.ID
				todo.Goal.Title = goal.Title
				newTodo["goal_id"] = goal.ID
				todo.Mission.ID = 0
				todo.Mission.Title = ""
				newTodo["mission_id"] = 0
			} else {
				todo.Goal.ID = 0
				todo.Goal.Title = ""
				newTodo["goal_id"] = 0
				todo.Mission.ID = 0
				todo.Mission.Title = ""
				newTodo["mission_id"] = 0
			}
		}
	}
	err = common.MySQLClient.
		Model(todo).
		Updates(newTodo).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
	todo, err := getTodoHelper(c.MustGet(common.ContextUserID), c.Param("id"))
	if err != nil {
		handleDBError(c, err, "todo")
		return
	}
	if err := common.MySQLClient.Delete(todo).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Status(http.StatusOK)
}
