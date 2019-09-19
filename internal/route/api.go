package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(router *gin.Engine)  {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello world")
	})
}
