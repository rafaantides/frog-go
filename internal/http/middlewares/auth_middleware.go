package middlewares

import (
	"frog-go/internal/core/domain"
	"frog-go/internal/core/errors"
	"frog-go/internal/utils/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(log *logger.Logger, validateToken func(string) (*domain.Claims, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("Authorization header missing")
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Message: errors.ErrorMessages[http.StatusUnauthorized],
				Detail:  "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Message: errors.ErrorMessages[http.StatusUnauthorized],
				Detail:  "Authorization header must be in the format 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := validateToken(token)
		if err != nil || claims == nil {
			log.Warn("Invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Message: errors.ErrorMessages[http.StatusUnauthorized],
				Detail:  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// adiciona o userID no contexto
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
