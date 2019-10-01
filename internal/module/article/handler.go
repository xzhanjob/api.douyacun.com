package article

import (
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListHandler(c *gin.Context) {
	total, data, err := NewIndex()
	if err != nil {
		logger.Errorf("首页文章列表错误: %s", err)
		c.JSON(http.StatusInternalServerError, "服务器出错了!")
		return
	}
	if total == 0 {
		c.JSON(http.StatusNotFound, gin.H{"total": total, "data": data})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "data": data})
	return
}

func InfoHandler(c *gin.Context) {
	id := c.Param("id")
	logger.Debugf("%s", id)
	data, err := NewInfo(id)
	if err != nil {
		logger.Errorf("%s", err)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if data == nil {
		c.JSON(http.StatusNotFound, "not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func TopicHandler(c *gin.Context) {
	topic := c.Param("topic")
	total, data, err := NewTopic(topic)
	if err != nil {
		logger.Errorf("首页文章列表错误: %s", err)
		c.JSON(http.StatusInternalServerError, "服务器出错了!")
		return
	}
	if total == 0 {
		c.JSON(http.StatusNotFound, gin.H{"total": total, "data": data})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "data": data})
	return
}
