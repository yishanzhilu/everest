package common

import (
	"github.com/gin-gonic/gin"
)

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
