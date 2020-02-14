package article

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"github.com/pkg/errors"
)

var Label _labels

type _labels struct {
	Label string `yaml:"label" json:"label"`
}

func (*_labels) List(size int) (data []string, err error) {
	var (
		buf bytes.Buffer
		r map[string]interface{}
	)
	query := map[string]interface{}{
		"_source": []string{"label"},
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"script": map[string]interface{}{
						"script": map[string]interface{}{
							"source": "doc['label'].size() > 0",
						},
					},
				},
			},
		},
	}
	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		panic(errors.Wrap(err, "json encode错误"))
	}

	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.TopicCost),
		db.ES.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if res.IsError() {
		panic(errors.Wrap(err, "es 查询错误"))
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		panic(errors.Wrap(err, "json decode错误"))
	}

	for _, v := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data = append(data, v.(map[string]interface{})["_source"].(map[string]interface{})["label"].(string))
	}
	return data, nil
}
