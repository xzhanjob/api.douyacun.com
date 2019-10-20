package article

import (
	"context"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"github.com/olivere/elastic/v7"
	"reflect"
	"time"
)

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

func NewIndex(p int) (int64, *[]index, error) {
	limit := 10
	skip := (p - 1) * limit
	logger.Debugf("skip: %v", skip)
	var (
		data index
		res  = make([]index, 0, 10)
	)
	fields := helper.GetStructJsonTag(data)
	_source := elastic.NewFetchSourceContext(true)
	_source.Include(fields...)
	searchResult, err := db.ES.Search().
		Index(TopicCost).
		FetchSourceContext(_source).
		From(skip).
		Size(limit).
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

func NewTopic(topic string) (int64, *[]index, error) {
	var (
		data index
		res  = make([]index, 0, 10)
	)
	fields := helper.GetStructJsonTag(data)
	tq := elastic.NewTermQuery("topic", topic)
	_source := elastic.NewSearchSource().
		Query(tq).
		FetchSource(true).
		FetchSourceIncludeExclude(fields, nil)
	searchResult, err := db.ES.Search().
		Index(TopicCost).
		SearchSource(_source).
		From(0).
		Size(10).
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
