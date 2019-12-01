package initialize

import (
	"dyc/internal/consts"
	"dyc/internal/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

var Config _config

type _config struct {
	MysqlDSN             string   `yaml:"mysql-dsn"`
	ElasticsearchAddress []string `yaml:"elstaticsearch-address"`
	Port                 string   `yaml:"port"`
	RunMode              string   `yaml:"run-mode"`
	Daemon               string   `yaml:"daemon"`
	PidFile              string   `yaml:"pid-file"`
	LogFile              string   `yaml:"log-file"`
	ImageDir             string   `yaml:"image-dir"`
}

var runMode = consts.DebugCode
var logFD *os.File

func (*_config) Init(filename string) *_config {
	var (
		err error
	)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatalf("Config file: %s", err)
	}
	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		logger.Fatal("Config file: %s", err)
	}
	if Config.IsRelease() {
		logFD, err = os.OpenFile(Config.Get().LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			logger.Fatalf("日志文件打开失败, %s", err)
		}
	} else {
		logFD = os.Stdout
	}
	return &Config
}

func (*_config) IsDaemon() bool {
	return strings.ToLower(Config.Daemon) == "on" || strings.ToLower(Config.Daemon) == "true"
}

func (*_config) IsRelease() bool {
	return runMode == consts.ReleaseCode
}

func (*_config) Get() *_config {
	return &Config
}

func (*_config) GetLogFD() *os.File {
	return logFD
}

func (*_config) SetRunMode(mode string) {
	switch mode {
	case consts.DebugMode:
		runMode = consts.DebugCode
	case consts.ReleaseMode:
		runMode = consts.ReleaseCode
	case consts.InfoMode:
		runMode = consts.TestCode
	}
}
