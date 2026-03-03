package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mmaruf23/go-task-management/internal/response"
)

func AuthMiddleware(j *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.AbortError(c, http.StatusUnauthorized, "authorization header required", nil)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.AbortError(c, http.StatusUnauthorized, "invalid authorization format", nil)
			return
		}

		tokenStr := parts[1]

		claims, err := j.VerifyToken(tokenStr)
		if err != nil {
			response.AbortError(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			response.AbortError(c, http.StatusUnauthorized, "invalid token claims", nil)
			return
		}
		// userID, err := uuid.Parse(claims.UserID) // todo : lanjutin yang ini. kerjain dulu jwtUtil nya
		c.Set("user_id", userID)

		c.Next()
	}
}
