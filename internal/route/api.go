package route

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine)  {
	api := router.Group("/api")
	{
		api.GET("/articles", )
	}
}
