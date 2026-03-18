package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mmaruf23/go-task-management/internal/response"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(s *AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Routes(r *gin.RouterGroup, authMiddlaware gin.HandlerFunc) {
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	protected := auth.Group("/", authMiddlaware)
	protected.PATCH("/password", h.UpdatePassword)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", response.ToErrorMap(err))

		return
	}

	token, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "REGISTER_SUCCESS", &token)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", response.ToErrorMap(err))

		return
	}

	token, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "LOGIN_SUCCESS", &token)
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}
	userID := value.(uuid.UUID)

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", response.ToErrorMap(err))
		return
	}

	err := h.service.UpdatePassword(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed update password", nil)
		return
	}

	response.Success[any](c, http.StatusOK, "RESET OK", nil)

}

// todo : implement refresh token
// - create new table "tokens"
// - create it's repository with sqlc
// - create it's service > create token, rotate, revoke, etc.
// - update register and login handler, should with refresh token as cookie.
// - create refresh, rotate, revoke token handler.
// - create unit test of each success scenario.

// do it besok aja. ya wkwk
