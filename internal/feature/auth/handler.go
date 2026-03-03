package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	// todo : do this thing
	// protected := auth.Group("/", authMiddlaware)
	// protected.GET("/me", h.Me)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}

// todo : implement refresh token
