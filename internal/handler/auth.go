package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/doc-storage/internal/config"
	"github.com/paudarco/doc-storage/internal/entity"
	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/paudarco/doc-storage/internal/handler/response"
	"github.com/paudarco/doc-storage/internal/service"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	userService service.User
	user        service.User
	cfg         *config.Config
	log         *logrus.Logger
}

func NewAuthHandler(userService service.User, user service.User, cfg *config.Config, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		user:        user,
		cfg:         cfg,
		log:         log,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	if req.AdminToken != h.cfg.AdminToken {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidAdminToken)
		return
	}

	if req.Login == "" || req.Pswd == "" {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	err := h.userService.Register(c.Request.Context(), req.Login, req.Pswd)
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"login": req.Login,
		},
	})
}

func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req entity.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	if req.Login == "" || req.Pswd == "" {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	token, err := h.userService.Authenticate(c.Request.Context(), req.Login, req.Pswd)
	if err != nil {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidCredentials)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"token": token,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.Param("token")

	err := h.user.InvalidateToken(c.Request.Context(), token)
	if err != nil {
		response.NewErrorResponse(c, h.log, fmt.Errorf("failed to invalidate token: %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			token: true,
		},
	})
}

func getUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", fmt.Errorf("user ID not found in context")
	}
	idStr, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("user ID is not a string")
	}
	return idStr, nil
}

func getQueryParams(c *gin.Context) (login, key, value string, limit int) {
	login = c.Query("login")
	key = c.Query("key")
	value = c.Query("value")
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	// Установим разумный лимит по умолчанию
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	return
}
