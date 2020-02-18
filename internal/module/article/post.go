package article

import (
	"bytes"
	"crypto/md5"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"time"
)

const PageSize = 10

var (
	Post _post
)

type _post struct{}

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

func (*_post) List(ctx *gin.Context, page int) (int64, []interface{}, error) {
	skip := (page - 1) * PageSize
	var (
		data  = make([]interface{}, 0, PageSize)
		total int64
		buf   bytes.Buffer
		r     map[string]interface{}
		err   error
	)
	query := map[string]interface{}{
		"from": skip,
		"size": PageSize,
		"sort": map[string]interface{}{
			"last_edit_time": map[string]interface{}{
				"order": "desc",
			},
		},
		"_source": []string{"author", "title", "description", "topic", "id", "cover", "date", "last_edit_time"},
	}
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode 错误"))
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if err != nil {
		return 0, data, errors.Wrap(consts.ESError{}, err.Error())
	}
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		panic(errors.Wrap(err, "json decode 错误"))
	}
	// 总条数
	total = int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, v := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data = append(data, v.(map[string]interface{})["_source"])
	}

	return total, data, nil
}

func (*_post) View(ctx *gin.Context, id string) (data interface{}, err error) {
	var (
		buf bytes.Buffer
		r   map[string]interface{}
	)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id.keyword": map[string]string{
					"value": id,
				},
			},
		},
	}
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode错误"))
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if err != nil {
		panic(errors.Wrap(err, "es请求错误"))
	}
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		panic(errors.Wrap(err, "json encode错误"))
	}

	data = r["hits"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["_source"]

	// 封面图片webp
	data.(map[string]interface{})["cover"] = Post.ConvertWebp(ctx, data.(map[string]interface{})["cover"].(string))
	// 内容图片webp
	data.(map[string]interface{})["content"] = Post.ConvertContentWebP(ctx, data.(map[string]interface{})["content"].(string))
	// 账户二维码webp
	data.(map[string]interface{})["wechat_subscription_qrcode"] = Post.ConvertWebp(ctx, data.(map[string]interface{})["wechat_subscription_qrcode"].(string))

	return
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
