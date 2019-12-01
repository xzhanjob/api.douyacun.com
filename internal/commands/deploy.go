package commands

import (
	"dyc/internal/consts"
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
			Name:     "conf",
			Usage:    "-conf [filename]",
			EnvVar:   "_DOUYACUN_CONF",
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
	initialize.Config.Init(c.String("conf"))
	// 设置运行环境
	initialize.Config.SetRunMode(initialize.Config.Get().RunMode)
	logger.NewLogger(initialize.Config.GetLogFD())
	logger.SetLevel(consts.DebugMode)
	// 数据库
	db.NewElasticsearch(initialize.Config.Get().ElasticsearchAddress)
	initialize.Config.SetRunMode(consts.DebugMode)
	deploy.Run(c.String("dir"))
	return
}
