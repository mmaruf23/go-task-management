package task

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description" binding:"required"`
}

type PaginationRequest struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

func (r *PaginationRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = 1
	}

	if r.Limit <= 0 {
		r.Limit = 10
	}

	if r.Limit > 100 {
		r.Limit = 100
	}
}

func (r *PaginationRequest) Offset() int32 {
	return (r.Limit) * (r.Page - 1)
}
