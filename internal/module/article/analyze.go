package article

import (
	"bufio"
	"bytes"
	"dyc/internal/config"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	TopicCost  = "articles"
	ImageRegex = `!\[(.*)\]\((.*)(.png|.gif|.jpg|.jpeg)(.*)\)`
)

// mapping
/**
{
    "mappings": {
        "properties": {
            "author": {
                "type": "keyword"
            },
            "content": {
                "type": "text"
            },
            "date": {
                "type": "date"
            },
            "description": {
                "type": "text"
            },
            "email": {
                "type": "keyword"
            },
            "github": {
                "type": "keyword",
                "index": false
            },
            "key": {
                "type": "keyword"
            },
            "keywords": {
                "type": "text"
            },
            "last_edit_time": {
                "type": "date"
            },
            "title": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "topic": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "wechat_subscription": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "wechat_subscription_qrcode": {
                "type": "text",
                "index": false
            }
        }
    }
}
*/

type article struct {
	Title                    string    `yaml:"Title" json:"title"`
	Keywords                 string    `yaml:"Keywords" json:"keywords"`
	Description              string    `yaml:"Description" json:"description"`
	Author                   string    `yaml:"Author" json:"author"`
	Date                     time.Time `yaml:"Date" json:"date"`
	LastEditTime             time.Time `yaml:"LastEditTime" json:"last_edit_time"`
	Content                  string    `yaml:"Content" json:"content"`
	Email                    string    `yaml:"Email" json:"email"`
	Github                   string    `yaml:"Github" json:"github"`
	Key                      string    `yaml:"Key" json:"key"`
	ID                       string    `yaml:"-" json:"-"`
	Topic                    string    `yaml:"-" json:"topic"`
	WechatSubscriptionQrcode string    `yaml:"WechatSubscriptionQrcode" json:"wechat_subscription_qrcode"`
	WechatSubscription       string    `yaml:"wechat_subscription" json:"wechat_subscription"`
}

func NewArticle(file string) (*article, error) {
	var (
		t   article
		err error
	)
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	conf := bytes.Buffer{}
	reader := bufio.NewReader(r)
	b, err := reader.ReadBytes('\n')
	if string(b[:len(b)-1]) == "---" {
		for {
			b, err = reader.ReadBytes('\n')
			if err != nil {
				return nil, err
			}
			if string(b[:len(b)-1]) == "---" {
				break
			}
			conf.Write(b)
		}
		//logger.Debugf("文章配置string: %s\n", conf.String())
		if err = yaml.Unmarshal(conf.Bytes(), &t); err != nil {
			return nil, err
		}
		//logger.Debugf("文章配置struct：%v", t)
	}
	var content = bytes.Buffer{}
	for {
		res, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		content.Write(res)
	}
	t.Content = content.String()
	return &t, nil
}

func (a *article) UploadImage(assert, dir string) (err error) {
	matched, err := regexp.MatchString(ImageRegex, a.Content)
	if err != nil {
		return errors.New(fmt.Sprintf("regexp match failed: %s", err))
	}
	if matched {
		// web访问图片目录
		imageDir := fmt.Sprintf("/%s/%s", a.Key, strings.Trim(dir, "/"))
		// 服务器存储目录
		storageDir := fmt.Sprintf("/%s%s", strings.Trim(config.Get().ImageDir, "/"), imageDir)
		if err = os.MkdirAll(storageDir, 0755); err != nil {
			return err
		}
		re, _ := regexp.Compile(ImageRegex)
		for _, v := range re.FindAllStringSubmatch(a.Content, -1) {
			filename := strings.Trim(v[2]+v[3], "/")
			src := fmt.Sprintf("%s/%s", assert, filename)
			// 替换文件image路径
			rebuild := strings.ReplaceAll(v[0], v[2]+v[3], fmt.Sprintf("%s/%s", imageDir, filename))
			// 服务器文件据对
			dst := storageDir + "/" + filename
			if !helper.FileExists(src) {
				logger.Warnf("image(%s) not found(%s)", v[0], src)
			}
			_, err = helper.Copy(dst, src)
			if err != nil {
				return err
			}
			logger.Debugf("上传图片 src: %s -> dst: %s", src, dst)
			a.Content = strings.ReplaceAll(a.Content, v[0], rebuild)
		}
	}
	return nil
}

func (a *article) Complete(c *conf, t string, position int) {
	if len(strings.TrimSpace(a.Author)) == 0 {
		a.Author = c.Author
	}
	if len(strings.TrimSpace(a.Github)) == 0 {
		a.Github = c.Github
	}
	if len(strings.TrimSpace(a.Email)) == 0 {
		a.Email = c.Email
	}
	a.WechatSubscription = c.WechatSubscription
	a.WechatSubscriptionQrcode = c.WechatSubscriptionQrcode
	a.Topic = strings.ToLower(t)
	a.Key = c.Key
	// todo
	// 1. git读取文章的创建时间和修改时间
	// 拼接文章id md5(user.id-topic-文章位置)
	a.ID =  fmt.Sprintf("%s-%s-%d", a.Topic, a.Key, position)
}

func (a *article) Storage() error {
	s, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	resp, err := db.ES.Index(
		TopicCost,
		strings.NewReader(string(s)),
		db.ES.Index.WithDocumentID(a.ID),
		db.ES.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		logger.Errorf("%s", data)
	}
	return nil
}
