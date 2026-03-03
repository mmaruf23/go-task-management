package response

import "github.com/gin-gonic/gin"

func Success[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(status, ApiResponse[T]{
		Success: true,
		Message: message,
		Data:    &data,
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
