package handler

import (
	"errors"
	"frog-go/internal/config"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/utils"
	"frog-go/internal/utils/utilsctx"
	"frog-go/internal/utils/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /api/v1/users/me
func (h *UserHandler) GetProfile(c *gin.Context) { ... }

// PUT /api/v1/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) { ... }

// PATCH /api/v1/users/me/password
func (h *UserHandler) UpdatePassword(c *gin.Context) { ... }

// PATCH /api/v1/users/me/email
func (h *UserHandler) UpdateEmail(c *gin.Context) { ... }

// DELETE /api/v1/users/me
func (h *UserHandler) DeactivateAccount(c *gin.Context) { ... }

// POST /api/v1/users/logout
func (h *UserHandler) Logout(c *gin.Context) { ... }