package initialize

import (
	"context"
	"dyc/internal/config"
	"dyc/internal/db"
	"dyc/internal/logger"
	"dyc/internal/route"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sync"
)

func Server(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()
	var (
		engine *gin.Engine
		err    error
	)
	// 日志
	fp, err := writer()
	if err != nil {
		log.Fatalf("log file writer: %s", err)
	}
	defer fp.Close()
	logger.NewLogger(fp)
	// 数据库
	db.NewDB(config.Get().MysqlDSN)
	defer shutdown()
	// 启动gin
	if config.IsRelease() {
		logger.SetLevel(logger.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.RecoveryWithWriter(fp))
	} else {
		logger.SetLevel(logger.InfoLevel)
		engine = gin.New()
		engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: logger.GinLogFormatter,
			Output:    fp,
			SkipPaths: nil,
		}), gin.RecoveryWithWriter(fp))
	}
	// 路由
	route.NewRouter(engine)
	server := http.Server{
		Addr:     config.Get().Port,
		Handler:  engine,
		ErrorLog: nil,
	}
	go func() {
		logger.Debug("start server 127.0.0.1:9001")
		_ = server.ListenAndServe()
	}()
	<-ctx.Done()
	if err := server.Close(); err != nil {
		logger.Fatalf("web server close failed: %s", err)
	}
}

func writer() (f *os.File, err error) {
	if config.IsRelease() {
		logFile := "/Users/liuning/Documents/github/douyacun-go/runtime/logs/douyacun.log"
		return os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		return os.Stdout, nil
	}
}
