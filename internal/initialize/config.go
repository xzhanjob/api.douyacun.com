package initialize

import (
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	ini "gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

var Config *ini.File
var logFd *os.File

func Init(env string) *ini.File {
	var (
		err error
	)
	switch env {
	case "debug":
		gin.SetMode(gin.DebugMode)
		//ini
		Config, err = ini.Load("configs/debug.ini")
		logger.SetLevel("debug")
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		Config, err = ini.Load("configs/prod.ini")
		logger.SetLevel("error")
	default:
		gin.SetMode(gin.ReleaseMode)
		Config, err = ini.Load("configs/prod.ini")
	}
	if err != nil {
		log.Fatalf("load config filed, %s", err)
	}
	logFd, err = os.OpenFile(GetKey("path::log_file").String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("path::log_file open failed, %s", err)
	}
	return Config
}

func GetKey(key string) *ini.Key {
	parts := strings.Split(key, "::")
	section := parts[0]
	keyStr := parts[1]
	return Config.Section(section).Key(keyStr)
}

func GetLogFD() *os.File  {
	return logFd
}