package middlewares

import (
	"frog-go/internal/core/errors"
	"frog-go/internal/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		handleError(c, log)
	}
}

func handleError(c *gin.Context, log *logger.Logger) {
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			log.Error("%v", err.Err)
		}

		err := c.Errors.Last().Err

		if appErr, ok := err.(*errors.AppError); ok {
			statusCode := appErr.StatusCode

			c.JSON(statusCode, errors.ErrorResponse{
				Message: appErr.Message,
				Detail:  appErr.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
				Message: errors.ErrorMessages[http.StatusInternalServerError],
				Detail:  err.Error(),
			})
		}
	}
}
