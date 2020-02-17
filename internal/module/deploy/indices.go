package deploy

import (
	"dyc/internal/consts"
	"dyc/internal/db"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var (
	Indices _Indices
)

type _article struct{}

type _Indices struct {
	Article _article
}

func (*_article) Create(index string) error {
	res, err := db.ES.Indices.Exists(
		[]string{index},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		res2, err := db.ES.Indices.Create(
			index,
			db.ES.Indices.Create.WithBody(strings.NewReader(consts.IndicesTopicMapping)),
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

func (*_article) Delete(index string) error {
	res, err := db.ES.Indices.Exists(
		[]string{index},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		res2, err := db.ES.Indices.Delete(
			[]string{index},
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

func (*_article) ReindexAndDeleteSource(source, dest string) (err error) {
	query := fmt.Sprintf(`{
	  "source": {
		"index": "%s"
	  },
	  "dest": {
		"index": "%s"
	  }
	}`, source, dest)
	res, err := db.ES.Reindex(
		strings.NewReader(query),
	)
	if err != nil {
		return errors.Wrap(err, "es reindex错误")
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(resp))
	}
	return Indices.Article.Delete(source)
}
