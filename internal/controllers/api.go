package controllers

import (
	"dyc/internal/middleware"
	"dyc/internal/module/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(router *gin.Engine) {
	hub := chat.NewHub()
	go hub.Run()
	api := router.Group("/api")
	{
		// 文章
		api.GET("/articles", Article.List)
		api.GET("/articles/labels", Article.Labels)
		api.GET("/article/:id", Article.View)
		api.GET("/topic/:topic", Topic.List)
		api.GET("/search/articles", Article.Search)
		api.POST("/subscribe", Subscribe.Create)
		// 电影资源
		api.GET("/media/subtype/:subtype", Media.Index)
		api.GET("/search/media", Media.Search)
		api.GET("/video/:id", Media.View)
		api.GET("/oauth/github", Oauth.Github)
		api.POST("/oauth/google", Oauth.Google)
		// websocket
		ws := api.Group("/ws", middleware.LoginCheck())
		{
			ws.GET("/join", func(context *gin.Context) {
				WS.Join(context, hub)
			})
			ws.POST("/channel", Channel.Private)
		}
	}
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
