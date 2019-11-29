package commands

import (
	"dyc/internal/initialize"
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
	initialize.Loading(c.String("conf"))
	deploy.Run(c.String("dir"))
	return
}
