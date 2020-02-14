package deploy

import (
	"dyc/internal/consts"
	"dyc/internal/db"
	"errors"
	"io/ioutil"
	"strings"
)

var (
	Topic _topic
)

type _topic struct{}

func (*_topic) Init() error {
	res, err := db.ES.Indices.Exists(
		[]string{consts.TopicCost},
	)
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		res2, err := db.ES.Indices.Create(
			consts.TopicCost,
			db.ES.Indices.Create.WithBody(strings.NewReader(consts.TopicMapping)),
		)
		if err != nil {
			return err
		}
		if res2.IsError() {
			resp, _ := ioutil.ReadAll(res.Body)
			return errors.New(string(resp))
		}
	}
	return nil
}
