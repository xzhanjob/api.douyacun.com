package route

import (
	"dyc/internal/module/article"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/articles", article.ListHandler)
		api.GET("/article/:id", article.InfoHandler)
		api.GET("/topic/:topic", article.TopicHandler)
		api.GET("/search", article.SearchHandler)
	}
}
