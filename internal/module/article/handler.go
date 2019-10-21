package article

import (
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func ListHandler(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := NewIndex(page)
	if err != nil {
		logger.Errorf("首页文章列表错误: %s", err)
		c.JSON(http.StatusInternalServerError, "服务器出错了!")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": total, "data": data}})
	return
}

func InfoHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := NewInfo(id)
	if err != nil {
		logger.Errorf("%s", err)
		c.JSON(http.StatusInternalServerError, "服务器出错了")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"data": data}})
}

func TopicHandler(c *gin.Context) {
	topic := c.Param("topic")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := NewTopic(topic, page)
	if err != nil {
		logger.Errorf("首页文章列表错误: %s", err)
		c.JSON(http.StatusInternalServerError, "服务器出错了!")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": total, "data": data}})
	return
}

func SearchHandler(c *gin.Context) {
	q := c.Query("q")
	if len(q) == 0 {
		c.JSON(http.StatusBadRequest, "请指定查询内容")
		return
	}
	total, data, err := NewSearch(q)
	if err != nil {
		logger.Errorf("文章搜索错误: %s", err)
		c.JSON(http.StatusInternalServerError, "服务器出错了!")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": total, "data": data}})
	return
}

func LabelHandler(c *gin.Context)  {
	// 关键字数量
	size := 30
	count := strings.TrimSpace(c.Param("count"))
	if count != "" {
		n, err := strconv.Atoi(count)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "非法请求"})
		}
		size = n
	}
	labels, _ := NewLabels(size)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": labels})
}