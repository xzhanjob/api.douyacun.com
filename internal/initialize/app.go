package initialize

import (
	"dyc/internal/config"
	"dyc/internal/db"
	"dyc/internal/logger"
	"fmt"
	"os"
)

func Loading(conf string)  {
	// 加载配置文件
	config.NewConfig(conf)
	// 设置运行环境
	config.SetRunMode(config.Get().RunMode)
	fp, err := writer()
	if err != nil {
		fmt.Printf("log file nil")
	}
	logger.NewLogger(fp)
	logger.SetLevel(config.Get().RunMode)
	// 数据库
	db.NewDB(config.Get().MysqlDSN)
	db.NewElasticsearch(config.Get().ElasticsearchAddress)
}

func writer() (f *os.File, err error) {
	if config.IsRelease() {
		return os.OpenFile(config.Get().LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		return os.Stdout, nil
	}
}
