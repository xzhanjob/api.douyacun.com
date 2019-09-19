package logger

import (
	"github.com/gin-gonic/gin"
	"time"
)

func GinLogFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	Infof("|%3d |%13v | %15s |%-7s %s %s",
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
	return ""
}
