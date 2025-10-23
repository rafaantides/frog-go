package handler

import (
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService inbound.AuthService
	userService inbound.UserService
}

func NewAuthHandler(authService inbound.AuthService, userService inbound.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

// Login godoc
// @Summary Login
// @Description Autentica o usuário pelo username **ou** email e senha, retornando um token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciais de login (identifier = username or email)"
// @Success 200 {object} dto.LoginResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	ctx := c.Request.Context()

	isEmail := func(s string) bool {
		_, err := mail.ParseAddress(s)
		return err == nil
	}

	var user *domain.User
	var err error

	if isEmail(req.Identifier) {
		user, err = h.userService.GetUserByEmail(ctx, req.Identifier)
	} else {
		user, err = h.userService.GetUserByUsername(ctx, req.Identifier)
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username/email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username/email or password"})
		return
	}

	token, err := h.authService.GenerateToken(ctx, user.ID, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.Error(appError.NewAppError(http.StatusBadRequest, err))
		return
	}

	data, err := h.userService.CreateUser(ctx, *input)
	if err != nil {
		c.Error(appError.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusCreated, data)
}
