package deploy

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"dyc/internal/config"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
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
	FilePath                 string    `yaml:"-" json:"-"`
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
	t.FilePath = file
	conf := bytes.Buffer{}
	reader := bufio.NewReader(r)
	b, err := reader.ReadBytes('\n')
	// 解析文件配置
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
func (a *Article) UploadImage(dir string) (err error) {
	// 文章图片所在目录绝对路径
	//assertAbs := fmt.Sprintf("/%s/%s", strings.Trim(dir, "/"), strings.Trim(assert, "/"))
	imageDir := helper.Path.Join(dir)
	// 图片服务存储目录
	//storageDir := fmt.Sprintf("/%s/%s/%s", strings.Trim(config.Get().ImageDir, "/"), a.Key, strings.Trim(assert, "/"))
	storageDir := helper.Path.Join(config.Get().ImageDir, a.Key)
	// 文章封面 -> 上传
	if len(a.Cover) > 0 {
		if err = os.MkdirAll(storageDir, 0755); err != nil {
			return
		}
		if _, err = helper.Copy(helper.Path.Join(storageDir, a.Cover), helper.Path.Join(imageDir, a.Cover)); err != nil {
			return
		}
		a.Cover = helper.Path.Join("images", a.Key, "assert", a.Cover)
		logger.Debugf("文章: %s 封面: %s", a.Title, a.Cover)
	} else {
		a.Cover = ""
	}
	// markdown图片 -> 上传
	matched, err := regexp.MatchString(consts.MarkDownImageRegex, a.Content)
	if err != nil {
		return errors.New(fmt.Sprintf("regexp match failed: %s", err))
	}
	if matched {
		if err = os.MkdirAll(storageDir, 0755); err != nil {
			return err
		}
		re, _ := regexp.Compile(consts.MarkDownImageRegex)
		for _, v := range re.FindAllStringSubmatch(a.Content, -1) {
			filename := strings.Trim(v[2]+v[3], "/")
			src := helper.Path.Join(imageDir, filename)
			// 替换文件image路径
			rebuild := strings.ReplaceAll(v[0], v[2]+v[3], helper.Path.Join("images", a.Key, filename))
			logger.Debugf("markdown 内部图片替换: %s -> %s", v[0], rebuild)
			// 服务器文件
			dst := helper.Path.Join(storageDir, filename)
			if !helper.FileExists(src) {
				logger.Warnf("image(%s) not found(%s)", v[0], src)
			}
			if _, err = helper.Copy(dst, src); err != nil {
				return
			}
			logger.Debugf("文章: %s 上传图片 src: %s -> dst: %s", a.Title, src, dst)
			a.Content = strings.ReplaceAll(a.Content, v[0], rebuild)
		}
	}
	return
}

// 完善信息文章
func (a *Article) Complete(c *Conf, topicTitle string, fileName string) {
	if strings.TrimSpace(a.Author) == "" {
		a.Author = c.Author
	}
	if strings.TrimSpace(a.Github) == "" {
		a.Github = c.Github
	}
	if strings.TrimSpace(a.Email) == "" {
		a.Email = c.Email
	}
	// 如果文章头部没有读取到标题，使用文件名作为标题
	if strings.TrimSpace(a.Title) == "" {
		a.Title = helper.Path.File(a.FilePath)
	}
	// 通过git版本获取最后更新时间
	lastEditTime, _ := helper.Git.LogFileLastCommitTime(a.FilePath)
	if a.LastEditTime.Before(lastEditTime) {
		a.LastEditTime = lastEditTime
	}
	// 通过git版本获取首次创建时间
	firstCreateTime, _ := helper.Git.LogFileFirstCommitTime(a.FilePath)
	if a.Date.After(firstCreateTime) {
		a.Date = firstCreateTime
	}
	// 每篇文章冗余一下二维码
	a.WechatSubscription = c.WechatSubscription
	a.WechatSubscriptionQrcode = c.WechatSubscriptionQrcode
	a.Topic = strings.ToLower(topicTitle)
	a.Key = c.Key
	// 1. git读取文章的创建时间和修改时间
	// 拼接文章id md5(user.key-topic-文件名称)
	a.ID = hex.EncodeToString(md5.New().Sum([]byte(fmt.Sprintf("%s-%s-%s", a.Topic, a.Key, fileName))))
}

// 存储文章
func (a *Article) Storage() error {
	s, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	_, err = db.ES.Index().
		Index(consts.TopicCost).
		BodyJson(string(s)).
		Id(a.ID).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
