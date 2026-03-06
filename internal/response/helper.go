package response

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mmaruf23/go-task-management/pkg/util"
)

func Success[T any, M util.PaginationMeta](c *gin.Context, status int, message string, data *T) {
	c.JSON(status, ApiResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
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

func Error(c *gin.Context, status int, message string, details *map[string][]string) {
	c.JSON(status, ApiResponse[any]{
		Success: false,
		Message: message,
		Errors:  details,
	})
}

func AbortError(c *gin.Context, status int, message string, details *map[string][]string) {
	Error(c, status, message, details)
	c.Abort()
}

func ToErrorMap(err error) *map[string][]string {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {

		errMap := make(map[string][]string)
		for _, f := range verr {
			field := strings.ToLower(f.Field())
			msg := fmt.Sprintf("failed on the '%s' tag", f.Tag())
			errMap[field] = append(errMap[field], msg)
		}
		return &errMap
	}

	return nil
}
