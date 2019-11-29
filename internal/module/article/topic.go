package article

import (
	"context"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"github.com/olivere/elastic/v7"
	"reflect"
)

var Topics _topic

type _topic struct{}

func (*_topic) List(topic string, page int) (int64, *[]index, error) {
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
		Index(consts.TopicCost).
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
