package commands

import (
	"bytes"
	"dyc/internal/config"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/logger"
	"encoding/csv"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"strings"
)

var AdCode = cli.Command{
	Name:        "adcode",
	Description: "获胜所属省级市级code",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "csv",
			Usage:    "--csv [path]",
			Required: true,
		},
		cli.StringFlag{
			Name:     "env",
			Usage:    "--env [env]",
			Required: false,
			Value:    "debug",
		},
	},
}

func adcode(c *cli.Context) error {
	config.Init(c.String("env"))
	db.NewElasticsearch(config.GetKey("elasticsearch::address").Strings(","), config.GetKey("elasticsearch::user").String(), config.GetKey("elasticsearch::password").String())
	fp, err := os.Open(c.String("csv"))
	if err != nil {
		return err
	}
	r := csv.NewReader(fp)
	if recodes, err := r.ReadAll(); err != nil {
		return err
	} else {

	}
}

func replaceIndices() {
	res, err := db.ES.Indices.Exists(
		[]string{consts.IndicesAdCodeConst},
	)
	if err != nil {
		logger.Errorf("index exists: %s", err)
		return
	}
	if res.StatusCode == 200 {
		res, err := db.ES.Indices.Delete(
			[]string{consts.IndicesAdCodeConst},
		)
		if err != nil {
			logger.Errorf("index delete: %s", err)
		}
		esResponsePrint(res)
	}
	if res, err := db.ES.Indices.Create(
		consts.IndicesAdCodeConst,
		db.ES.Indices.Create.WithBody(strings.NewReader(consts.IndicesAdCodeMapping)),
	); err != nil {
		logger.Errorf()
	}
}

func esResponsePrint(response *esapi.Response) {
	if response.IsError() {
		body, _ := ioutil.ReadAll(response.Body)
		logger.Errorf("es response error: %s", body)
	}
}

func storage(records [][]string) {
	header := records[0]
	start := 0
	name, adCode, cityCode := 0, 0, 0
	for k, v := range header {
		switch v {
		case "adcode":
			adCode = k
			start = 1
		case "citycode":
			cityCode = k
		case "中文名":
			name = k
		}
	}
	buf := bytes.NewBuffer(nil)
	indexBulk := fmt.Sprintf(`{ "index":{"_index":"%s"}}`, consts.IndicesAdCodeConst)
	for i := start; i < len(records); i++ {
		buf.WriteString(indexBulk)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf(`{"name": "%s", "adcode": %s, "citycode": %s}`, records[i][name], records[i][adCode], records[i][cityCode]))
		buf.WriteString("\n")
	}
	res, err := db.ES.Bulk(buf)
	if err != nil {
		logger.Error(err)
		return
	}
	esResponsePrint(res)
}
