package commands

import (
	"dyc/internal/db"
	"dyc/internal/initialize"
	"dyc/internal/logger"
	"dyc/internal/module/deploy"
	"github.com/urfave/cli"
)

var Deploy = cli.Command{
	Name: "deploy",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "env",
			Usage:    "-env [dev, prod]",
			EnvVar:   "_DOUYACUN_ENV",
			Required: true,
		},
		cli.StringFlag{
			Name:     "dir",
			Usage:    "指定文章所在目录",
			Required: true,
		},
	},
	Action: deployAction,
}

func deployAction(c *cli.Context) (err error) {
	// 加载配置文件
	initialize.Init(c.String("env"))
	// 设置运行环境
	logger.NewLogger(initialize.GetLogFD())
	// 数据库
	db.NewElasticsearch(initialize.GetKey("elasticsearch::address").Strings(","), initialize.GetKey("elasticsearch::user").String(), initialize.GetKey("elasticsearch::password").String())
	deploy.Run(c.String("dir"))
	return
}
