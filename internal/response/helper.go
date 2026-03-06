package response

import (
	"github.com/gin-gonic/gin"
	"github.com/mmaruf23/go-task-management/pkg/util"
)

func Success[T any, M util.PaginationMeta](c *gin.Context, status int, message string, data T) {
	c.JSON(status, ApiResponse[T]{
		Success: true,
		Message: message,
		Data:    &data,
	})
}

func SuccessWithMeta[T any](c *gin.Context, status int, message string, data T, meta util.PaginationMeta) {
	c.JSON(status, ApiResponse[T]{
		Success: true,
		Message: message,
		Data:    &data,
		Meta:    &meta,
	})
}

func Error(c *gin.Context, status int, message string, details any) {
	c.JSON(status, ApiResponse[any]{
		Success: false,
		Message: message,
		Error:   &details,
	})
}

func AbortError(c *gin.Context, status int, message string, details any) {
	Error(c, status, message, details)
	c.Abort()
}
