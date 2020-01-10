package workspace

import "github.com/gin-gonic/gin"

// TaskHandler .
type TaskHandler interface {
	GetTask(c *gin.Context)
	GetTaskList(c *gin.Context)
	PostTask(c *gin.Context)
	PutTask(c *gin.Context)
	PatchTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

// NewTaskHandler .
func NewTaskHandler() TaskHandler {
	return &taskHandler{}
}

// taskHandler .
type taskHandler struct {
}

func (h *taskHandler) GetTask(c *gin.Context) {

}

func (h *taskHandler) GetTaskList(c *gin.Context) {

}

func (h *taskHandler) PostTask(c *gin.Context) {

}

func (h *taskHandler) PutTask(c *gin.Context) {

}

func (h *taskHandler) PatchTask(c *gin.Context) {

}

func (h *taskHandler) DeleteTask(c *gin.Context) {

}
