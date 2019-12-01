package deploy

import (
	"dyc/internal/helper"
	"dyc/internal/initialize"
	"dyc/internal/logger"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Root 文章所在目录
// Author 文章作者
// Email 作者邮箱
// Github 作者github连接
// WechatSubscription 微信订阅号
// WechatSubscriptionQrcode 微信订阅号二维码
// Key 唯一标识符
type Conf struct {
	Root                     string              `yaml:"-"`
	Author                   string              `yaml:"Author"`
	Email                    string              `yaml:"Email"`
	Github                   string              `yaml:"Github"`
	WechatSubscription       string              `yaml:"WechatSubscription"`
	WechatSubscriptionQrcode string              `yaml:"WechatSubscriptionQrcode"`
	Key                      string              `yaml:"Key"`
	Topics                   map[string][]string `yaml:"Topics"`
}

// 加载配置文件
func LoadConfig(dir string) (*Conf, error) {
	var (
		conf Conf
	)
	configFile := path.Join(dir, "douyacun.yml")
	conf.Root = dir
	logger.Debugf("配置文件路径: %s", configFile)
	if !helper.File.IsFile(configFile) {
		return nil, errors.New("请先配置douyaucn.yaml")
	}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(b, &conf); err != nil {
		return nil, err
	}
	logger.Debugf("%v", conf)
	return &conf, nil
}

//
func (c *Conf) UploadQrcode(dir string) (err error) {
	// 服务器存储目录
	storageDir := fmt.Sprintf("/%s/%s/%s", strings.Trim(initialize.Config.Get().ImageDir, "/"), c.Key, filepath.Dir(c.WechatSubscriptionQrcode))
	if err = os.MkdirAll(storageDir, 0755); err != nil {
		return err
	}
	src := fmt.Sprintf("/%s/%s", strings.Trim(dir, "/"), strings.Trim(c.WechatSubscriptionQrcode, "/"))
	dst := fmt.Sprintf("/%s/%s/%s", strings.Trim(initialize.Config.Get().ImageDir, "/"), c.Key, strings.Trim(c.WechatSubscriptionQrcode, "/"))
	logger.Debugf("上传二维码 src: %s -> dst: %s", src, dst)
	_, err = helper.File.Copy(dst, src)
	if err != nil {
		return err
	}
	c.WechatSubscriptionQrcode = fmt.Sprintf("/%s/%s/%s", "images", c.Key, strings.Trim(c.WechatSubscriptionQrcode, "/"))
	return nil
}
