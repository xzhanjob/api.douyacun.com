package article

import (
	"context"
	"crypto/md5"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const PageSize = 10

var (
	Post _post
)

type _post struct{}

type index struct {
	Author       string    `json:"author"`
	Date         time.Time `json:"date"`
	LastEditTime time.Time `json:"last_edit_time"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Topic        string    `json:"topic"`
	Id           string    `json:"id"`
	Cover        string    `json:"cover"`
}

type view struct {
	Title                    string    `json:"title"`
	Keywords                 string    `json:"keywords"`
	Label                    string    `json:"label"`
	Cover                    string    `json:"cover"`
	Description              string    `json:"description"`
	Author                   string    `json:"author"`
	Date                     time.Time `json:"date"`
	LastEditTime             time.Time `json:"last_edit_time"`
	Content                  string    `json:"content"`
	Email                    string    `json:"email"`
	Github                   string    `json:"github"`
	Key                      string    `json:"key"`
	ID                       string    `json:"id"`
	Topic                    string    `json:"topic"`
	FilePath                 string    `json:"-"`
	WechatSubscriptionQrcode string    `json:"wechat_subscription_qrcode"`
	WechatSubscription       string    `json:"wechat_subscription"`
}

func (*_post) List(ctx *gin.Context, page int) (int64, []index, error) {
	skip := (page - 1) * PageSize
	var (
		data index
		res  = make([]index, 0, PageSize)
	)
	fields := helper.GetStructJsonTag(data)
	_source := elastic.NewFetchSourceContext(true)
	_source.Include(fields...)
	searchResult, err := db.ES.Search().
		Index(consts.TopicCost).
		FetchSourceContext(_source).
		Sort("last_edit_time", false).
		From(skip).
		Size(PageSize).
		Do(context.Background())
	if err != nil {
		return 0, nil, err
	}
	for k, item := range searchResult.Each(reflect.TypeOf(data)) {
		tmp := item.(index)
		tmp.Id = searchResult.Hits.Hits[k].Id
		tmp.Cover = Post.ConvertWebp(ctx, tmp.Cover)
		res = append(res, tmp)
	}
	return searchResult.TotalHits(), res, nil
}

func (*_post) View(ctx *gin.Context, id string) (*view, error) {
	resp, err := db.ES.Get().Index(consts.TopicCost).Id(id).Do(context.Background())
	if err != nil {
		return nil, nil
	}
	var data view
	if err = json.Unmarshal(resp.Source, &data); err != nil {
		return nil, err
	}
	if resp.Source == nil {
		return nil, nil
	}
	data.Cover = Post.ConvertWebp(ctx, data.Cover)
	data.WechatSubscriptionQrcode = Post.ConvertWebp(ctx, data.WechatSubscriptionQrcode)
	data.Content = Post.ConvertContentWebP(ctx, data.Content)
	logger.Debugf("%s", resp.Id)
	return &data, nil
}

func (*_post) ConvertWebp(ctx *gin.Context, image string) string {
	ext := path.Ext(image)
	if helper.Image.WebPSupportExt(ext) {
		ua := ctx.Request.UserAgent()
		if strings.Contains(ua, "Chrome") || strings.Contains(ua, "Android") {
			return strings.Replace(image, ext, ".webp", 1)
		}
	}
	return image
}

func (c *_post) ConvertContentWebP(ctx *gin.Context, content string) string {
	matched, err := regexp.MatchString(consts.MarkDownImageRegex, content)
	if err != nil {
		return content
	}
	if matched {
		re, _ := regexp.Compile(consts.MarkDownImageRegex)
		for _, v := range re.FindAllStringSubmatch(content, -1) {
			filename := v[2] + v[3]
			WebP := Post.ConvertWebp(ctx, filename)
			if WebP != filename {
				// 替换文件image路径
				rebuild := strings.ReplaceAll(v[0], v[2]+v[3], WebP)
				content = strings.ReplaceAll(content, v[0], rebuild)
			}
		}
	}
	return content
}

// 拼接文章id md5(user.key-topic-文件名称)
func (c *_post) GenerateId(topic, key, filename string) string {
	return hex.EncodeToString(md5.New().Sum([]byte(fmt.Sprintf("%s-%s-%s", topic, key, filename))))
}
