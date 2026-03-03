package response

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
	Error   *any   `json:"error,omitempty"`
}
