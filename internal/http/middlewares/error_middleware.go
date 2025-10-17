package middlewares

import (
	appError "frog-go/internal/core/errors"
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

		if appErr, ok := err.(*appError.AppError); ok {
			statusCode := appErr.StatusCode

			c.JSON(statusCode, appError.ErrorResponse{
				Message: appErr.Message,
				Detail:  appErr.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, appError.ErrorResponse{
				Message: appError.ErrorMessages[http.StatusInternalServerError],
				Detail:  err.Error(),
			})
		}
	}
}
