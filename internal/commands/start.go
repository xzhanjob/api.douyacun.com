package commands

import (
	"context"
	"dyc/internal/config"
	"dyc/internal/helper"
	"dyc/internal/initialize"
	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var Start = cli.Command{
	Name:   "start",
	Usage:  "",
	Action: startAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "conf",
			Usage:    "-conf [filename]",
			EnvVar:   "_DOUYACUN_CONF",
			Required: true,
		},
	},
}

func startAction(c *cli.Context) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 加载配置文件
	config.NewConfig(c.String("conf"))
	// 设置运行环境
	config.SetRunMode(config.Get().RunMode)

	dmn := &daemon.Context{
		PidFileName: config.Get().PidFile,
		PidFilePerm: 0644,
		LogFileName: "",
		LogFilePerm: 0,
		WorkDir:     "/",
		Umask:       027,
	}
	if !daemon.WasReborn() && config.IsDaemon() {
		cancel()
		if pid, ok := initialize.ChildAlreadyRunning(config.Get().PidFile); ok {
			log.Fatalf("daemon already running with process id %v", pid)
		}
		child, err := dmn.Reborn()
		if err != nil {
			log.Fatal(err)
		}
		if child != nil {
			if !helper.FileOverwrite(config.Get().PidFile, []byte(strconv.Itoa(child.Pid))) {
				log.Fatalf("failed writing process id to \"%s\"", config.Get().PidFile)
			}
			log.Fatalf("daemon started with process id %v\n", child.Pid)
		}
	}

	// 启动web服务
	go initialize.Server(ctx, &wg)
	// 处理结束信号
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Print("web server shutdown...")
	// 通知goroutine要结束，关闭一下资源
	cancel()
	wg.Done()
	return nil
}
