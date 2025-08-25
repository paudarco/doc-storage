package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/paudarco/doc-storage/internal/config"
	"github.com/paudarco/doc-storage/internal/service"
	"github.com/sirupsen/logrus"
)

type Auth interface {
	Authenticate(c *gin.Context)
	Register(c *gin.Context)
	Logout(c *gin.Context)
}

type Doc interface {
	UploadDoc(c *gin.Context)
	ListDocs(c *gin.Context)
	GetDoc(c *gin.Context)
	DeleteDoc(c *gin.Context)
}

type Handler struct {
	Doc
	Auth
}

func NewHandler(service *service.Service, cfg *config.Config, log *logrus.Logger) *Handler {
	return &Handler{
		Doc:  NewDocHandler(service.Doc, log),
		Auth: NewAuthHandler(service.User, service.User, cfg, log),
	}
}
