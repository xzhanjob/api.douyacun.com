package article

import (
	"context"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"github.com/olivere/elastic/v7"
	"reflect"
)

var Label _labels

type _labels struct {
	Label string `yaml:"label" json:"label"`
}

func (*_labels) List(size int) (l *[]string, err error) {
	var (
		data _labels
		res  = make([]string, 0, size)
	)
	q := elastic.NewBoolQuery().Filter(elastic.NewScriptQuery(elastic.NewScript(`doc['label'].size() > 0`)))
	_source := elastic.NewSearchSource().Query(q).
		FetchSource(true).
		FetchSourceIncludeExclude(helper.GetStructJsonTag(data), nil)
	searchResult, err := db.ES.Search().
		Index(consts.TopicCost).
		SearchSource(_source).
		From(0).
		Size(size).
		Do(context.Background())
	logger.Debugf("%s", searchResult.Hits.Hits[0].Source)
	if err != nil {
		return nil, err
	}
	for _, item := range searchResult.Each(reflect.TypeOf(data)) {
		res = append(res, item.(_labels).Label)
	}
	return &res, nil
}
