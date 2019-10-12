package config

import (
	"dyc/internal/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Conf struct {
	MysqlDSN             string   `yaml:"mysql-dsn"`
	ElasticsearchAddress []string `yaml:"elstaticsearch-address"`
	Port                 string   `yaml:"port"`
	RunMode              string   `yaml:"run-mode"`
	Daemon               string   `yaml:"daemon"`
	PidFile              string   `yaml:"pid-file"`
	LogFile              string   `yaml:"log-file"`
	ImageDir             string   `yaml:"image-dir"`
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

func IsDaemon() bool {
	if strings.ToLower(config.Daemon) == "on" || strings.ToLower(config.Daemon) == "true" {
		return true
	} else {
		return false
	}
}

func Get() *Conf {
	return &config
}
