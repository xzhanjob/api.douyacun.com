package subscribe

import (
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type formEmail struct {
	Email string `json:"email"`
}

func CreateHandler(c *gin.Context) {
	var (
		err error
		e   formEmail
	)
	if err = c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404})
		return
	}
	logger.Debugf("%v", e)
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", e.Email); !m {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "邮件格式错误"})
		return
	}
	if err = NewSubscriber(e.Email).Store(); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "存储失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "成功"})
}
