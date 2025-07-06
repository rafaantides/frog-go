package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"frog-go/internal/core/errors"
	"frog-go/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UUIDMiddleware(lg *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, values := range c.Request.URL.Query() {
			if strings.HasSuffix(strings.ToLower(key), "id") {
				for _, value := range values {
					if _, err := uuid.Parse(value); err != nil {
						abortWithError(c, lg, http.StatusBadRequest, errors.InvalidParam(key, err))
						return
					}
				}
			}
		}

		if err := validateBodyUUIDs(c); err != nil {
			abortWithError(c, lg, http.StatusBadRequest, err)
			return
		}

		c.Next()
	}
}

func validateBodyUUIDs(c *gin.Context) error {
	if c.Request.ContentLength == 0 || !strings.HasPrefix(c.ContentType(), "application/json") {
		return nil
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		return err
	}

	return validateUUIDsRecursive(body)
}

func validateUUIDsRecursive(data map[string]interface{}) error {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			if strings.HasSuffix(strings.ToLower(key), "id") {
				if _, err := uuid.Parse(v); err != nil {
					return errors.InvalidParam(key, err)
				}
			}
		case []interface{}:
			for _, item := range v {
				switch el := item.(type) {
				case string:
					if _, err := uuid.Parse(el); err != nil {
						return errors.InvalidParam(key, err)
					}
				case map[string]interface{}:
					if err := validateUUIDsRecursive(el); err != nil {
						return err
					}
				}
			}
		case map[string]interface{}:
			if err := validateUUIDsRecursive(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func abortWithError(c *gin.Context, lg *logger.Logger, status int, err error) {
	c.Error(errors.NewAppError(status, err))
	handleError(c, lg)
	c.Abort()
}
