package initialize

import (
	"context"
	"dyc/internal/controllers"
	"dyc/internal/db"
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Server(ctx context.Context) {
	var (
		engine *gin.Engine
	)
	// 日志
	logger.NewLogger(Config.GetLogFD())
	// 数据库
	db.NewElasticsearch(Config.Get().ElasticsearchAddress)
	defer shutdown()
	// 启动gin
	if Config.IsRelease() {
		logger.SetLevel(Config.Get().RunMode)
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.RecoveryWithWriter(Config.GetLogFD()))
	} else {
		logger.SetLevel(Config.Get().RunMode)
		engine = gin.New()
		engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: logger.GinLogFormatter,
			Output:    Config.GetLogFD(),
			SkipPaths: nil,
		}), gin.RecoveryWithWriter(Config.GetLogFD()))
	}
	// 路由
	controllers.NewRouter(engine)
	server := http.Server{
		Addr:     Config.Get().Port,
		Handler:  engine,
		ErrorLog: nil,
	}
	go func() {
		logger.Debugf("start server 127.0.0.1%s", Config.Get().Port)
		_ = server.ListenAndServe()
	}()
	<-ctx.Done()
	if err := server.Close(); err != nil {
		logger.Fatalf("web server close failed: %s", err)
	}
}
