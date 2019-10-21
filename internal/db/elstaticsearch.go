package db

import (
	"github.com/olivere/elastic/v7"
	"sync"
)

var ES *elastic.Client
var es_once sync.Once

func NewElasticsearch(address []string) {
	var err error
	es_once.Do(func() {
		ES, err = elastic.NewClient(elastic.SetURL(address...), elastic.SetSniff(false))
		if err != nil {
			panic(err)
		}
	})
}

func NewElasticsearchDefault()  {
	var err error
	es_once.Do(func() {
		ES, err = elastic.NewClient(elastic.SetURL("http://192.168.1.2:9200"), elastic.SetSniff(false))
		if err != nil {
			panic(err)
		}
	})
}