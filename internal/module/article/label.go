package article

import (
	"context"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"github.com/olivere/elastic/v7"
	"reflect"
)

type Labels struct {
	Label string `yaml:"label" json:"label"`
}

func NewLabels(size int) (l *[]string, err error) {
	var (
		data Labels
		res  = make([]string, 0, size)
	)
	q := elastic.NewBoolQuery().Filter(elastic.NewScriptQuery(elastic.NewScript(`doc['label'].size() > 0`)))
	_source := elastic.NewSearchSource().Query(q).
		FetchSource(true).
		FetchSourceIncludeExclude(helper.GetStructJsonTag(data), nil)
	searchResult, err := db.ES.Search().
		Index(TopicCost).
		SearchSource(_source).
		From(0).
		Size(size).
		Do(context.Background())
	logger.Debugf("%s", searchResult.Hits.Hits[0].Source)
	if err != nil {
		return nil, err
	}
	for _, item := range searchResult.Each(reflect.TypeOf(data)) {
		res = append(res, item.(Labels).Label)
	}
	return &res, nil
}
