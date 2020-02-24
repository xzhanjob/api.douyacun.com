package initialize

import (
	"bytes"
	"context"
	"dyc/internal/controllers"
	"dyc/internal/db"
	"dyc/internal/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Server() {
	var (
		engine *gin.Engine
		err error
	)
	// 日志
	logger.NewLogger(GetLogFD())
	// 数据库
	db.NewElasticsearch(GetKey("elasticsearch::address").Strings(","), GetKey("elasticsearch::user").String(), GetKey("elasticsearch::password").String())
	defer shutdown()
	// 启动gin
	engine = gin.New()
	//engine.Use(recoverWithWrite(GetLogFD()))
	engine.Use(gin.RecoveryWithWriter(GetLogFD()))

	// 路由
	controllers.NewRouter(engine)
	port := GetKey("server::port").String()
	server := http.Server{
		Addr:     port,
		Handler:  engine,
		ErrorLog: nil,
	}
	go func() {
		logger.Infof("start server 127.0.0.1%s", port)
		if err = server.ListenAndServe(); err != nil {
			logger.Info("web server closed")
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("web server close failed: %s", err)
	}
	time.Sleep(50 * time.Millisecond)
}

func recoverWithWrite(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				if brokenPipe {
					logger.Debugf("%s\n%s%s", err, string(httpRequest), logger.Reset)
				} else {
					buf := new(bytes.Buffer) // the returned data
					err := errors.WithStack(err.(error))
					fmt.Fprintf(buf,"%+v", err)
					logger.Errorf("[Recovery] panic recovered:\n%s\n%s", strings.Join(headers, "\r\n"),  buf.String())
				}
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器出错了"})
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
