package article

import (
	"context"
	"dyc/internal/db"
	"dyc/internal/helper"
	"github.com/olivere/elastic/v7"
	"reflect"
	"time"
)
const PageSize = 10

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

func NewIndex(page int) (int64, *[]index, error) {
	skip := (page - 1) * PageSize
	var (
		data index
		res  = make([]index, 0, PageSize)
	)
	fields := helper.GetStructJsonTag(data)
	_source := elastic.NewFetchSourceContext(true)
	_source.Include(fields...)
	searchResult, err := db.ES.Search().
		Index(TopicCost).
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

func NewTopic(topic string, page int) (int64, *[]index, error) {
	var (
		data index
		res  = make([]index, 0, PageSize)
	)
	skip := (page - 1) * PageSize
	fields := helper.GetStructJsonTag(data)
	tq := elastic.NewTermQuery("topic", topic)
	_source := elastic.NewSearchSource().
		Query(tq).
		FetchSource(true).
		FetchSourceIncludeExclude(fields, nil).
		Sort("last_edit_time", false)
	searchResult, err := db.ES.Search().
		Index(TopicCost).
		SearchSource(_source).
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
