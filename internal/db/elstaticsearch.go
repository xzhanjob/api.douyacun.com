package db

import (
	"dyc/internal/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"sync"
)

var ES *elasticsearch.Client
var es_once sync.Once

func NewElasticsearch(address []string) {
	var err error
	es_once.Do(func() {
		ES, err = elasticsearch.NewClient(elasticsearch.Config{Addresses: address})
		if err != nil {
			logger.Fatal("elstaticsearch new client failed: %s", err)
		}
	})
}
