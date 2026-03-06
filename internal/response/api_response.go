package response

import "github.com/mmaruf23/go-task-management/pkg/util"

type ApiResponse[T any] struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    *T                   `json:"data,omitempty"`
	Meta    *util.PaginationMeta `json:"meta,omitempty"`
	Error   *any                 `json:"error,omitempty"`
}
