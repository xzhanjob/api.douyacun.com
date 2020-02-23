package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/articles", Article.List)
		api.GET("/articles/labels", Article.Labels)
		api.GET("/article/:id", Article.View)
		api.GET("/topic/:topic", Topic.List)
		api.GET("/search/articles", Article.Search)
		api.POST("/subscribe", Subscribe.Create)
		api.GET("/media/subtype/:subtype", Media.Index)
		api.GET("/search/media", Media.Search)
		api.GET("/video/:id", Media.View)
		api.GET("/panic/test", Media.Test)
	}
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
