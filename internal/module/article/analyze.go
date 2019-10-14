package article

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"dyc/internal/config"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	TopicCost    = "articles"
	ImageRegex   = `!\[(.*)\]\((.*)(.png|.gif|.jpg|.jpeg)(.*)\)`
	TopicMapping = `{
    "mappings": {
        "properties": {
            "author": {
                "type": "keyword"
            },
            "content": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "date": {
                "type": "date"
            },
            "description": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
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
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "label": {
                "type": "text",
                "index": true,
				"fielddata": true
            },
            "last_edit_time": {
                "type": "date"
            },
            "title": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart",
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
}`
)

type Article struct {
	Title                    string    `yaml:"Title" json:"title"`
	Keywords                 string    `yaml:"Keywords" json:"keywords"`
	Label                    string    `yaml:"Label" json:"label"`
	Cover                    string    `yaml:"Cover" json:"cover"`
	Description              string    `yaml:"Description" json:"description"`
	Author                   string    `yaml:"Author" json:"author"`
	Date                     time.Time `yaml:"Date" json:"date"`
	LastEditTime             time.Time `yaml:"LastEditTime" json:"last_edit_time"`
	Content                  string    `yaml:"Content" json:"content"`
	Email                    string    `yaml:"Email" json:"email"`
	Github                   string    `yaml:"Github" json:"github"`
	Key                      string    `yaml:"Key" json:"key"`
	ID                       string    `yaml:"-" json:"id"`
	Topic                    string    `yaml:"-" json:"topic"`
	WechatSubscriptionQrcode string    `yaml:"WechatSubscriptionQrcode" json:"wechat_subscription_qrcode"`
	WechatSubscription       string    `yaml:"wechat_subscription" json:"wechat_subscription"`
}

func NewArticle(file string) (*Article, error) {
	var (
		t   Article
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

// assert: douyacun.yml 配置的话题图片存储目录 dir: book所在目录
func (a *Article) UploadImage(dir, assert string) (err error) {
	// 文章图片所在目录绝对路径
	assertAbs := fmt.Sprintf("/%s/%s", strings.Trim(dir, "/"), strings.Trim(assert, "/"))
	// 图片服务存储目录
	storageDir := fmt.Sprintf("/%s/%s/%s", strings.Trim(config.Get().ImageDir, "/"), a.Key, strings.Trim(assert, "/"))
	if len(a.Cover) > 0 {
		// 文章封面
		_, err = helper.Copy(storageDir+"/"+a.Cover, assertAbs+"/"+a.Cover)
		if err != nil {
			return err
		}
		a.Cover = fmt.Sprintf("/%s/%s/%s/%s", "images", a.Key, strings.Trim(assert, "/"), a.Cover)
		logger.Debugf("文章: %s 封面: %s", a.Title, a.Cover)
	} else {
		a.Cover = ""
	}
	// markdown图片
	matched, err := regexp.MatchString(ImageRegex, a.Content)
	if err != nil {
		return errors.New(fmt.Sprintf("regexp match failed: %s", err))
	}
	if matched {

		if err = os.MkdirAll(storageDir, 0755); err != nil {
			return err
		}
		re, _ := regexp.Compile(ImageRegex)
		for _, v := range re.FindAllStringSubmatch(a.Content, -1) {
			filename := strings.Trim(v[2]+v[3], "/")
			src := fmt.Sprintf("%s/%s", assertAbs, filename)
			// 替换文件image路径
			rebuild := strings.ReplaceAll(v[0], v[2]+v[3], fmt.Sprintf("/%s/%s/%s/%s", "images", a.Key, strings.Trim(assert, "/"), filename))
			logger.Debugf("markdown image replace: %s -> %s", v[0], rebuild)
			// 服务器文件
			dst := storageDir + "/" + filename
			if !helper.FileExists(src) {
				logger.Warnf("image(%s) not found(%s)", v[0], src)
			}
			_, err = helper.Copy(dst, src)
			if err != nil {
				return err
			}
			logger.Debugf("文章: %s 上传图片 src: %s -> dst: %s", a.Title, src, dst)
			a.Content = strings.ReplaceAll(a.Content, v[0], rebuild)
		}
	}
	return nil
}

func (a *Article) Complete(c *conf, t string, position int) {
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
	a.ID = hex.EncodeToString(md5.New().Sum([]byte(fmt.Sprintf("%s-%s-%d", a.Topic, a.Key, position))))
}

func (a *Article) Storage() error {
	s, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	_, err = db.ES.Index().
		Index(TopicCost).
		BodyJson(string(s)).
		Id(a.ID).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func Initialize() error {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := db.ES.IndexExists(TopicCost).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		_, err := db.ES.CreateIndex(TopicCost).Body(TopicMapping).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func Purge(key string) error {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := db.ES.IndexExists(TopicCost).Do(context.Background())
	if err != nil {
		return err
	}
	if exists {
		bq := elastic.NewTermQuery("key", key)
		purgeResp, err := db.ES.DeleteByQuery().Index(TopicCost).Query(bq).Do(context.Background())
		if err != nil {
			return err
		}
		logger.Debugf("delete articles: %v", purgeResp)
	}
	return nil
}
