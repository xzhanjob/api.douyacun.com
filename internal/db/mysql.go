package db

import (
	"dyc/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var write *gorm.DB

func InitMysql(dsn string) {
	// videos_t
	var err error
	write, err = gorm.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("gorm open %s", err)
	}
	if err = write.DB().Ping(); err != nil {
		logger.Fatal(err)
	}
	write.SingularTable(true)
}

func Write(ctx *gin.Context) *gorm.DB {
	return write
}

func Close() {
	_ = write.Close()
}
