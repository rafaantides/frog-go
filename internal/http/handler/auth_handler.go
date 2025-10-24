package handler

import (
	"context"
	"fmt"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/ports/inbound"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	token, err := h.createTokenAndSession(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, *token)
}

func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input, err := req.ToDomain()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(ctx, *input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.createTokenAndSession(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, *token)
}

func (h *AuthHandler) createTokenAndSession(ctx context.Context, userID uuid.UUID) (*string, error) {
	duration := time.Hour * 1

	token, err := h.authService.GenerateToken(ctx, userID, duration)
	if err != nil {
		return nil, fmt.Errorf("could not generate token")
	}

	session, err := domain.NewUserSession(userID, token)
	if err != nil {
		return nil, fmt.Errorf("could not generate user session")
	}

	if err := h.authService.CreateUserSession(ctx, *session); err != nil {
		return nil, fmt.Errorf("could not create user session")
	}

	return &token, nil
}
