package common

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	timeFormart = "2006-01-02T15:04:05Z"
)

// BaseModel extends gorm.Model, but change ID to sortable uuid
type BaseModel struct {
	ID        string     `json:"id" gorm:"PRIMARY_KEY;type:char(20);not null;column:id;comment:'主键id'"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// BaseSerialize extends gorm.Model, but change ID to sortable uuid
type BaseSerialize struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// MarshalJSON ..
func (m *BaseModel) MarshalJSON() ([]byte, error) {
	type Alias BaseModel
	return json.Marshal(&BaseSerialize{
		ID:        m.ID,
		CreatedAt: m.CreatedAt.Format(timeFormart),
		UpdatedAt: m.UpdatedAt.Format(timeFormart),
	})
}

// BaseHandler can register gin router group
type BaseHandler interface {
	RegisterPublicRoutes(r *gin.RouterGroup)
	RegisterPrivateRoutes(r *gin.RouterGroup)
}

// ResponseError contains error message, description, and statusCode
type ResponseError struct {
	message     string
	description interface{}
	statusCode  int
}

// NewResponseError create new ResponseError with specific info
func NewResponseError(message string, description interface{}, statusCode int) *ResponseError {
	return &ResponseError{
		message,
		description,
		statusCode,
	}
}

// NewResponseErrorWithErr create new ResponseError, with err as description
func NewResponseErrorWithErr(message string, err error, statusCode int) *ResponseError {
	return &ResponseError{
		message,
		err,
		statusCode,
	}
}

func (r ResponseError) Error() string {
	return r.message
}

// Abort will abort passed gin reqest and respose the error
func (r ResponseError) Abort(c *gin.Context) {
	c.AbortWithStatusJSON(r.statusCode,
		gin.H{"error": r.message, "error_description": r.description},
	)
}
