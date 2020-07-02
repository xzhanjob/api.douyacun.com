package deploy

import (
	"bufio"
	"bytes"
	"dyc/internal/config"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"dyc/internal/module/article"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
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

// bookDir: book所在目录
// book图片存储目录 {topic}
// 服务器图片存储目录 /images/blog/{a.Key}/{topic}/{image}
// 配置图片存储目录 /data/web/images/blog/
// /images/blog/ 需要作为前缀
// 注意: 这里image topic为根目录，一般是 assert/a.jpg
func (a *Article) UploadImage(bookDir string, topic string) (err error) {
	// 图片前缀
	imagePrefix := path.Join("/images/blog", a.Key, topic)
	// 图片服务存储目录, 去掉images，方便后面直接拼接images
	storageDir := config.GetKey("path::storage_dir").String()
	var errTemplate = func(s string) error {
		return fmt.Errorf("《%s》: %s", a.Title, s)
	}
	var debugTemplate = func(s string) {
		logger.Debugf("《%s》: %s", a.Title, s)
	}
	var warnTemplate = func(s string) {
		logger.Warnf("《%s》: %s", a.Title, s)
	}
	// 文章封面 -> 上传
	if len(a.Cover) > 0 {
		if _, err = helper.File.Copy(path.Join(storageDir, imagePrefix, a.Cover), path.Join(bookDir, topic, a.Cover)); err != nil {
			return fmt.Errorf("article %s 封面复制失败, %s", a.Title, err)
		}
		a.Cover = path.Join(imagePrefix, a.Cover)
		debugTemplate(fmt.Sprintf("封面复制成功: %s", a.Cover))
	} else {
		a.Cover = ""
	}
	// markdown图片 -> 上传
	matched, err := regexp.MatchString(consts.MarkDownImageRegex, a.Content)
	if err != nil {
		return errTemplate(fmt.Sprintf("regexp match failed: %s", err))
	}
	if matched {
		re, err := regexp.Compile(consts.MarkDownImageRegex)
		if err != nil {
			return errTemplate(fmt.Sprintf("regex compile faile: %s", err))
		}
		for _, v := range re.FindAllStringSubmatch(a.Content, -1) {
			filename := strings.Trim(v[2]+v[3], "/")
			src := path.Join(bookDir, topic, filename)
			// 替换文件image路径
			rebuild := strings.ReplaceAll(v[0], v[2]+v[3], path.Join(imagePrefix, filename))
			debugTemplate(fmt.Sprintf("image: %s -> %s", v[0], rebuild))
			// 服务器文件
			dst := path.Join(storageDir, imagePrefix, filename)
			if !helper.File.IsFile(src) {
				warnTemplate(fmt.Sprintf("image %s not found(%s)", v[0], src))
				continue
			}
			if _, err = helper.File.Copy(dst, src); err != nil {
				return errTemplate(fmt.Sprintf("image copy failed, %s", err))
			}
			debugTemplate(fmt.Sprintf("image upload src: %s -> dst: %s", src, dst))
			a.Content = strings.ReplaceAll(a.Content, v[0], rebuild)
		}
	}
	// 本地文件跳转
	mathchLocal, err := regexp.MatchString(consts.MarkDownLocalJump, a.Content)
	if err != nil {
		return errTemplate(fmt.Sprintf("regex(%s) match string failed: %s", consts.MarkDownLocalJump, err))
	}
	if mathchLocal {
		re, err := regexp.Compile(consts.MarkDownLocalJump)
		if err != nil {
			return errTemplate(fmt.Sprintf("regex(%s) compile failed: %s", consts.MarkDownLocalJump, err))
		}
		var targetTopic string
		for _, v := range re.FindAllStringSubmatch(a.Content, -1) {
			// 目标文件
			target := v[1]
			// 话题 以顶层目录为话题
			if path.IsAbs(target) { // 绝对路径
				// 判断目标文件是否存在
				if !helper.File.IsFile(path.Join(bookDir, target)) {
					warnTemplate(fmt.Sprintf("本地跳转(%s)不存在", path.Join(bookDir, target)))
					continue
				}
				targetTopic = helper.File.TopDir(target)
				if targetTopic == "" {
					warnTemplate(fmt.Sprintf("本地跳转，目标文件: %s, 一级目录不存在，暂不支持跳转", target))
					continue
				}
				if targetTopic != strings.Trim(path.Dir(target), "/") {
					warnTemplate(fmt.Sprintf("本地跳转: %s，一层目录是话题，文章必须在同一目录下，暂不支持多层目录", target))
					continue
				}
			} else { // 相对路径
				// 判断目标文件是否存在
				if !helper.File.IsFile(path.Join(bookDir, a.Topic, target)) {
					warnTemplate(fmt.Sprintf("本地跳转(%s) 文件不存在", path.Join(bookDir, a.Topic, target)))
					continue
				}
				if path.Dir(target) != "." {
					warnTemplate(fmt.Sprintf("本地跳转: %s，一层目录是话题，文章必须在同一目录下，暂不支持多层目录", target))
				}
				targetTopic = a.Topic
			}
			targetFileName := path.Base(target)
			targetID := article.Post.GenerateId(targetTopic, a.Key, targetFileName)
			targetUrl := fmt.Sprintf("/article/%s", targetID)
			replaceContent := strings.ReplaceAll(v[0], target, targetUrl)
			debugTemplate(fmt.Sprintf("本地跳转：%s -> %s", v[0], replaceContent))
			a.Content = strings.ReplaceAll(a.Content, v[0], replaceContent)
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
		a.Title = path.Base(a.FilePath)
	}
	var cst, _ = time.LoadLocation("Asia/Shanghai")
	// 通过git版本获取最后更新时间
	lastEditTime, _ := helper.Git.LogFileLastCommitTime(a.FilePath)
	if a.LastEditTime.Before(lastEditTime) {
		a.LastEditTime = lastEditTime
	}
	z, _ := a.LastEditTime.Zone()
	if z != cst.String() {
		a.LastEditTime = a.LastEditTime.Add(-time.Hour * 8)
	}
	logger.Debugf("《%s》: 更新时间: %s, 时区: %s", a.Title, a.LastEditTime, z)
	// 通过git版本获取首次创建时间
	firstCreateTime, _ := helper.Git.LogFileFirstCommitTime(a.FilePath)
	if a.Date.After(firstCreateTime) {
		a.Date = firstCreateTime
	}
	z, _ = a.Date.Zone()
	if z != cst.String() {
		a.Date = a.Date.Add(-time.Hour * 8)
	}
	logger.Debugf("《%s》: 创建时间: %s, 时区: %s", a.Title, a.LastEditTime, z)
	// 每篇文章冗余一下二维码
	a.WechatSubscription = c.WechatSubscription
	a.WechatSubscriptionQrcode = c.WechatSubscriptionQrcode
	a.Topic = strings.ToLower(topicTitle)
	a.Key = c.Key
	a.ID = article.Post.GenerateId(a.Topic, a.Key, fileName)
}

// 存储文章
func (a *Article) Storage(index string) (err error) {
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(a); err != nil {
		return
	}
	res, err := db.ES.Index(
		index,
		strings.NewReader(buf.String()),
	)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(resp))
	}
	logger.Debugf("storage: 《%s》存储成功", a.Title)
	return
}
