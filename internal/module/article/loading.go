package article

import (
	"dyc/internal/config"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
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

func (c *conf) Qrcode(dir string) (err error) {
	// 服务器存储目录
	storageDir := fmt.Sprintf("/%s/%s/%s", strings.Trim(config.Get().ImageDir, "/"), c.Key, filepath.Dir(c.WechatSubscriptionQrcode))
	if err = os.MkdirAll(storageDir, 0755); err != nil {
		return err
	}
	src := fmt.Sprintf("/%s/%s", strings.Trim(dir, "/"), strings.Trim(c.WechatSubscriptionQrcode, "/"))
	dst := fmt.Sprintf("/%s/%s/%s", strings.Trim(config.Get().ImageDir, "/"), c.Key, strings.Trim(c.WechatSubscriptionQrcode, "/"))
	logger.Debugf("上传二维码 src: %s -> dst: %s", src, dst)
	_, err = helper.Copy(dst, src)
	if err != nil {
		return err
	}
	c.WechatSubscriptionQrcode = fmt.Sprintf("/%s/%s/%s", "images", c.Key, strings.Trim(c.WechatSubscriptionQrcode, "/"))
	return nil
}
