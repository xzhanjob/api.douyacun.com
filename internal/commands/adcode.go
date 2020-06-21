package commands

import (
	"bytes"
	"dyc/internal/config"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/csv"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
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
	db.ES.Indices.Exists(
		[]string{consts.IndicesAdCodeConst},
	)
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
	for _, record := range records {
		buf.WriteString(indexBulk)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf(`{"name": "%s", "adcode": %s, "citycode": %s}`, record[name], record[adCode], record[cityCode]))
		buf.WriteString("\n")
	}
	res, err := db.ES.Bulk(buf)
	if err != nil {
		log.Error(err)
		return
	}
	if res.IsError() {
		msg, _ := ioutil.ReadAll(res.Body)
		log.Error("es bulk create error: %s", msg)
		return
	}
}
