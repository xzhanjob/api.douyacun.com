package config

import (
	"dyc/internal/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Conf struct {
	MysqlDSN string `yaml:"mysql-dsn"`
	Port     string `yaml:"port"`
	RunMode  string `yaml:"run-mode"`
	Daemon   string `yaml:"daemon"`
	PidFile  string `yaml:"pid-file"`
}

var config Conf

func NewConfig(filename string) *Conf {
	var (
		err error
	)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatalf("config file: %s", err)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		logger.Fatal("config file: %s", err)
	}
	return &config
}

func IsDaemon() bool  {
	if strings.ToLower(config.Daemon) == "on" {
		return true
	} else {
		return false
	}
}

func Get() *Conf {
	return &config
}
