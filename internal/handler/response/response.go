package response

import (
	"github.com/gin-gonic/gin"
	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, log *logrus.Logger, err error) {
	errCode := errors.CheckError(err)
	message := err.Error()
	log.Error(message)
	c.AbortWithStatusJSON(errCode, gin.H{
		"error": errorResponse{errCode, message},
	})
}
