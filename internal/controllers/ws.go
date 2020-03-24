package controllers

import (
	"dyc/internal/module/chat"
	"github.com/gin-gonic/gin"
)

var WS _ws

type _ws struct{}


func (*_ws) Join(ctx *gin.Context) {
	chat.ServeWs(ctx)
}