package article

import (
	"dyc/internal/helper"
	"dyc/internal/logger"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Topic struct {
	Icon     string   `yaml:"Icon"`
	Dir      string   `yaml:"Dir"`
	Assert   string   `yaml:"Assert"`
	Articles []string `yaml:"Articles"`
}

type conf struct {
	Root                     string           `yaml:"-"`
	Author                   string           `yaml:"Author"`
	Email                    string           `yaml:"Email"`
	Github                   string           `yaml:"Github"`
	WechatSubscription       string           `yaml:"WechatSubscription"`
	WechatSubscriptionQrcode string           `yaml:"WechatSubscriptionQrcode"`
	Key                      string           `yaml:"Key"`
	Topics                   map[string]Topic `yaml:"Topics"`
}

func LoadDir(dir string) (*conf, error) {
	var (
		err error
		c   conf
	)
	dir = strings.TrimRight(dir, "/")
	c.Root = dir
	f := dir + "/douyacun.yml"
	if !helper.FileExists(f) {
		return nil, errors.New("请先配置douyaucn.yaml")
	}
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	logger.Debugf("%v", c)
	return &c, nil
}
