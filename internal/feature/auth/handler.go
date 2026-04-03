package auth

import (
	"fmt"
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
	auth.POST("/refresh", h.Refresh)

	protected := auth.Group("/", authMiddlaware)
	protected.PATCH("/password", h.UpdatePassword)
	protected.POST("/logout", h.Logout)
	protected.POST("/logout-all", h.LogoutAll)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", response.ToErrorMap(err))
		return
	}

	userID, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	newJti := uuid.New()
	token, err := h.service.GenerateToken(c.Request.Context(), newJti, uuid.MustParse(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error generate token", nil)
		return
	}

	c.SetCookie("refresh_token", token.Refresh, token.MaxAgeRefereshToken, "/", "", false, true) // note : enable https kalau untuk prod.

	response.Success(c, http.StatusOK, "REGISTER_SUCCESS", &token.Access)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", response.ToErrorMap(err))

		return
	}

	userID, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	newJti := uuid.New()
	token, err := h.service.GenerateToken(c.Request.Context(), newJti, uuid.MustParse(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error generate token", nil)
		return
	}

	c.SetCookie("refresh_token", token.Refresh, token.MaxAgeRefereshToken, "/", "", false, true) // note : enable https kalau untuk prod.

	response.Success(c, http.StatusOK, "LOGIN_SUCCESS", &token.Access)
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

// todo : handle juga kalau tokennya udah dipake. > cek replaced / revoked.
// next : akhirnya rombak besar-besaran untuk yang ini. nanti rewrite juga yang generateToken service. supaya nggak bikin id nya oleh database. untuk save ke database, nanti dibikin aja service baru. untuk generate khusus generate token. alurnya nanti : validasi token lama, gen jti, rotate db, gen token, done.
func (h *AuthHandler) Refresh(c *gin.Context) {

	oldToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "no token provided", nil)
		return
	}

	claims, err := h.service.jwt.VerifyToken(oldToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid token format", nil)
		return
	}
	userID := uuid.MustParse(claims.Subject)

	newJti := uuid.New()
	err = h.service.ReplaceToken(c.Request.Context(), uuid.MustParse(claims.ID), newJti)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	token, err := h.service.GenerateToken(c.Request.Context(), newJti, userID)
	if err != nil {
		fmt.Printf("CODE II : %s", err.Error())
		response.Error(c, http.StatusInternalServerError, "error refresh token : CODE II", nil)
		return
	}

	c.SetCookie("refresh_token", token.Refresh, token.MaxAgeRefereshToken, "/", "", false, true) // note : enable https kalau untuk prod.
	response.Success(c, http.StatusOK, "SUCCESS_REFRESH_TOKEN", &token.Access)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "no token provided", nil)
		return
	}

	err = h.service.Logout(c.Request.Context(), token)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie("refresh_token", "", 0, "/", "", false, true) // note : gak tau ini bener atau salah, kalau logout cookienya diginiin?
	response.Success[any](c, http.StatusOK, "LOGOUT_SUCCESS", nil)
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}
	userID := value.(uuid.UUID)

	err := h.service.LogoutAll(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.SetCookie("refresh_token", "", 0, "/", "", false, true) // note : gak tau ini bener atau salah, kalau logout cookienya diginiin?
	response.Success[any](c, http.StatusOK, "LOGOUT_ALL_DEVICE_SUCCESS", nil)
}

// ngantuk,
