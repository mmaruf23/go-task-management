package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mmaruf23/go-task-management/internal/response"
)

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// ROUTE
func (h *TaskHandler) Routes(r *gin.RouterGroup, authMiddlaware gin.HandlerFunc) {
	task := r.Group("/task", authMiddlaware)

	task.POST("/", h.Create)
	task.GET("/", h.List)
	task.PATCH("/:id", h.Status)
}

// HANDLER METHOD
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		// todo : detailnya error validationnya masukin
		return
	}

	value, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}
	userID := value.(uuid.UUID)

	task, err := h.service.CreateTask(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Success create new task", task)
}

func (h *TaskHandler) List(c *gin.Context) {
	var req PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Normalize()

	value, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}
	userID := value.(uuid.UUID)

	results, err := h.service.GetUserTasks(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Success get user task list", results.Data, results.Meta)
}

func (h *TaskHandler) Status(c *gin.Context) {
	var req TaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	status, err := req.Parse()
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	value, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}
	userID := value.(uuid.UUID)

	taskID, err := uuid.Parse(c.Param("id")) // note : yang kaya gini mau juga kah dibikin helper? *edit : ya, sekalian validasi format id nya => uuid
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	if err = h.service.UpdateStatus(c.Request.Context(), userID, taskID, status); err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success[any](c, http.StatusOK, "success update status", nil)

}
