package deploy

import (
	"context"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/logger"
	"github.com/olivere/elastic/v7"
)

var (
	Topic _topic
)

type _topic struct{}

func (*_topic) Init() error {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := db.ES.IndexExists(consts.TopicCost).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		_, err := db.ES.CreateIndex(consts.TopicCost).Body(consts.TopicMapping).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (*_topic) Purge(key string) error {
	// Use the IndexExists service to check if a specified index exists.
	exists, err := db.ES.IndexExists(consts.TopicCost).Do(context.Background())
	if err != nil {
		return err
	}
	if exists {
		bq := elastic.NewTermQuery("key", key)
		purgeResp, err := db.ES.DeleteByQuery().Index(consts.TopicCost).Query(bq).Do(context.Background())
		if err != nil {
			return err
		}
		logger.Debugf("delete articles: %v", purgeResp)
	}
	return nil
}
