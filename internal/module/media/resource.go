package media

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var Resource _Resource

type _Resource struct{}

func (*_Resource) Index(page int, subtype string) (total int64, data []interface{}, err error) {
	skip := (page - 1) * consts.MediaDefaultPageSize
	var (
		buf bytes.Buffer
		r   map[string]interface{}
	)
	query := map[string]interface{}{
		"from": skip,
		"size": consts.MediaDefaultPageSize,
		"sort": map[string]interface{}{
			"updated_at": map[string]interface{}{
				"order": "desc",
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"torrent_num": map[string]interface{}{"gte": 0},
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"subtype": subtype,
						},
					},
				},
			},
		},
		"_source": []string{"id", "title", "region", "genres", "released", "rate", "summary"},
	}
	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		panic(errors.Wrap(err, "json encode 错误"))
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMediaConst),
		db.ES.Search.WithBody(&buf),
	)
	if err != nil {
		panic(errors.Wrap(err, "es查询错误"))
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		panic(errors.Wrap(err, "json decode 错误"))
	}
	total = int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	for _, v := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data = append(data, Resource.ToMap(v.(map[string]interface{})["_source"]))
	}
	return
}

func (*_Resource) ToMap(data interface{}) interface{} {
	if v, ok := data.(map[string]interface{}); ok {
		var (
			f = map[string]interface{}{
				"id":          v["id"],
				"rate":        int64(v["rate"].(float64)),
				"title":       v["title"],
				"released":    v["released"],
				"description": v["summary"],
				"author":      "",
				"cover":       "",
			}
			genres []string
		)
		if f, ok := v["genres"].([]interface{}); ok {
			if len(f) > 2 {
				f = f[:2]
			}
			if len(f) == 0 {
				if f, ok = v["region"].([]interface{}); ok {
					if len(f) > 2 {
						f = f[:2]
					}
				}
			}
			for _, g := range f {
				genres = append(genres, g.(string))
			}
		}
		f["author"] = strings.Join(genres, "/")
		return f
	} else {
		return data
	}
}
