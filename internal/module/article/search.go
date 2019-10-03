package article

import (
	"context"
	"dyc/internal/db"
	"dyc/internal/helper"
	"github.com/olivere/elastic/v7"
	"reflect"
	"time"
)

type search struct {
	Author       string    `json:"author"`
	Date         time.Time `json:"date"`
	LastEditTime time.Time `json:"last_edit_time"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Topic        string    `json:"topic"`
	Id           string    `json:"id"`
	Highlight    []string  `json:"highlight"`
}

func NewSearch(q string) (int64, *[]search, error) {
	var (
		data search
		res  = make([]search, 0, 10)
	)
	_source := elastic.NewSearchSource().
		Highlight(elastic.NewHighlight().Field("content")).
		Query(elastic.NewMultiMatchQuery(q, "title", "content", "author", "keywords")).
		FetchSource(true).
		FetchSourceIncludeExclude(helper.GetStructJsonTag(data), nil)
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
		tmp := item.(search)
		tmp.Id = searchResult.Hits.Hits[k].Id
		tmp.Highlight = searchResult.Hits.Hits[k].Highlight["content"]
		res = append(res, tmp)
	}
	return searchResult.TotalHits(), &res, nil
}
