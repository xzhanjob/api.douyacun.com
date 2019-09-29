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
	db.NewElasticsearch(config.Get().ElasticsearchAddress)
	defer shutdown()
	// 启动gin
	if config.IsRelease() {
		logger.SetLevel(config.Get().RunMode)
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.RecoveryWithWriter(fp))
	} else {
		logger.SetLevel(config.Get().RunMode)
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
