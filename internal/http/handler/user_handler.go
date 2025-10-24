package handler

import (
	"errors"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/utils/utilsctx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service inbound.UserService
}

func NewUserHandler(service inbound.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GET /api/v1/users/me
func (h *UserHandler) GetProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusUnauthorized, err))
	}

	data, err := h.service.GetUser(ctx, userID)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, data)
}

// PUT /api/v1/users/me
func (h *UserHandler) UpdateProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusUnauthorized, err))
		return
	}

	var req struct {
		Name     string `json:"name" binding:"omitempty"`
		Username string `json:"username" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	updatedUser, err := h.service.UpdateUserProfile(ctx, userID, req.Name, req.Username)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// PATCH /api/v1/users/me/password
func (h *UserHandler) UpdatePasswordHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusUnauthorized, err))
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	err = h.service.UpdateUserPassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, appError.ErrInvalidPassword):
			c.Error(appError.NewAppError(http.StatusBadRequest, err))
		default:
			c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// PATCH /api/v1/users/me/email
func (h *UserHandler) UpdateEmailHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusUnauthorized, err))
		return
	}

	var req struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	err = h.service.UpdateUserEmail(ctx, userID, req.NewEmail)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email updated successfully",
	})
}

// DELETE /api/v1/users/me
func (h *UserHandler) DeactivateAccountHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusUnauthorized, err))
		return
	}

	err = h.service.DeactivateUserAccount(ctx, userID)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deactivated successfully",
	})
}

// POST /api/v1/users/logout
func (h *UserHandler) LogoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusUnauthorized, err))
		return
	}

	err = h.service.LogoutUser(ctx, userID)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
}
