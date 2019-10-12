package commands

import (
	"dyc/internal/config"
	"dyc/internal/logger"
	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli"
	"os"
	"syscall"
)

var StopCommand = cli.Command{
	Name:   "stop",
	Action: StopAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "conf",
			Usage:    "-conf [filename]",
			EnvVar:   "_DOUYACUN_CONF",
			Required: true,
		},
	},
}

func StopAction(c *cli.Context) error {
	// 加载配置文件
	config.NewConfig(c.String("conf"))
	logger.NewLogger(os.Stdout)
	logger.Infof("looking for pid in (%s)", config.Get().PidFile)
	dctx := new(daemon.Context)
	dctx.PidFileName = config.Get().PidFile
	child, err := dctx.Search()
	if err != nil {
		logger.Info(err)
		return nil
	}
	if err = child.Signal(syscall.SIGTERM); err != nil {
		logger.Info(err)
		return nil
	}
	st, err := child.Wait()
	if err != nil {
		logger.Info("daemon exited successfully")
		return nil
	}
	logger.Infof("daemon[%v] exited[%v]? successfully[%v]?\n", st.Pid(), st.Exited(), st.Success())
	return nil
}
