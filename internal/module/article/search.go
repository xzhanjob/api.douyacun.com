package article

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

var Search _search

type _search struct {
	Author       string    `json:"author"`
	Date         time.Time `json:"date"`
	LastEditTime time.Time `json:"last_edit_time"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Topic        string    `json:"topic"`
	Id           string    `json:"id"`
	Highlight    []string  `json:"highlight"`
}

func (*_search) List(q string) (total int64, data []interface{}, err error) {
	type response struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source struct {
					Date         time.Time `json:"date"`
					LastEditTime time.Time `json:"last_edit_time"`
					Author       string    `json:"author"`
					Topic        string    `json:"topic"`
					Id           string    `json:"id"`
					Title        string    `json:"title"`
				} `json:"_source"`
				Highlight struct {
					Content []string `json:"content"`
				} `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	var (
		buf bytes.Buffer
		r   response
	)

	data = make([]interface{}, 0)
	query := map[string]interface{}{
		"_source": []string{"author", "title", "description", "topic", "id", "date", "last_edit_time"},
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q,
				"fields": []string{"title.keyword", "author", "keywords", "content", "label"},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
		},
	}
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode 错误"))
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if err != nil {
		panic(errors.Wrap(err, "es search错误"))
	}
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		panic(errors.Wrap(err, "json decode 错误"))
	}
	total = r.Hits.Total.Value
	for _, v := range r.Hits.Hits {
		tmp := map[string]interface{}{
			"date":      v.Source.Date,
			"id":        v.Source.Id,
			"author":    v.Source.Author,
			"topic":     v.Source.Topic,
			"title":     v.Source.Title,
			"highlight": v.Highlight.Content,
		}
		data = append(data, tmp)
	}
	return
}
