package util

type Paginated[T any] struct {
	Data T              `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int32 `json:"total_pages"`
}

func BuildPaginationMeta(page, limit int32, total int64) *PaginationMeta {
	totalPages := int32((total + int64(limit) - 1) / int64(limit))

	return &PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
