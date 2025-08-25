package handler

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		auth := api.Group("/")
		{
			auth.POST("/register", h.Register)
			auth.POST("/auth", h.Authenticate)
		}

		authorized := api.Group("/")
		authorized.Use()

		authorized.DELETE("/auth/:token", h.Logout)

		docs := authorized.Group("/docs")
		{
			docs.POST("/", h.UploadDoc)
			docs.GET("/", h.ListDocs)
			docs.HEAD("/", h.ListDocs)
			docs.GET("/:id", h.GetDoc)
			docs.HEAD("/:id", h.GetDoc)
			docs.DELETE("/:id", h.DeleteDoc)
		}

	}

	return router
}
