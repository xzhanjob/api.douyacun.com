package route

import (
	"dyc/internal/module/article"
	"dyc/internal/module/subscribe"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/articles", article.ListHandler)
		api.GET("/articles/labels", article.LabelHandler)
		api.GET("/article/:id", article.InfoHandler)
		api.GET("/topic/:topic", article.TopicHandler)
		api.GET("/search", article.SearchHandler)
		api.POST("/subscribe", subscribe.CreateHandler)
	}
}
