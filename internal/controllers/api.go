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
	router.Use(middleware.Trace())
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
		api.GET("/media/info", Media.Info)
		// websocket
		auth := api.Group("/", middleware.LoginCheck())
		{
			auth.GET("/ws/join", func(context *gin.Context) {
				WS.Join(context, hub)
			})
			auth.POST("/ws/channel", Channel.Create)
			auth.GET("/ws/channel/subscribe", Channel.subscribe)
			auth.GET("/ws/channel/messages", Channel.Messages)
			auth.GET("/account/list", Account.List)
		}
	}
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
