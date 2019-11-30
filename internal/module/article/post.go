package article

import (
	"context"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"dyc/internal/module/deploy"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"path"
	"reflect"
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
	Cover        cover    `json:"cover"`
}

func (*_post) List(ctx *gin.Context, page int) (int64, *[]index, error) {
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
		res = append(res, tmp)
	}
	return searchResult.TotalHits(), &res, nil
}

func (*_post) View(id string) (*deploy.Article, error) {
	resp, err := db.ES.Get().Index(consts.TopicCost).Id(id).Do(context.Background())
	if err != nil {
		return nil, nil
	}
	var data deploy.Article
	if err = json.Unmarshal(resp.Source, &data); err != nil {
		return nil, err
	}
	if resp.Source == nil {
		return nil, nil
	}
	logger.Debugf("%s", resp.Id)
	return &data, nil
}

type cover string

func (c *cover) Convert(ctx *gin.Context) {
	ext := path.Ext(string(*c))
	if helper.Image.WebPSupportExt(ext) {
		*c = cover(strings.Replace(string(*c), ext, ".webp", 1))
	}
}
