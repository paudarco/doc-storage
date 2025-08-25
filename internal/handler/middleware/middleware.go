package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/paudarco/doc-storage/internal/handler/response"
	"github.com/paudarco/doc-storage/internal/service"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(userService service.UserService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Bearer realm="api"`)
			response.NewErrorResponse(c, log, errors.ErrTokenRequired)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.NewErrorResponse(c, log, errors.ErrInvalidAuthHeader)
			return
		}

		token := parts[1]

		userID, err := userService.ValidateToken(c, token)
		if err == errors.ErrTokenExpired {
			response.NewErrorResponse(c, log, err)
			return
		} else if err != nil {
			response.NewErrorResponse(c, log, errors.ErrInvalidToken)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
