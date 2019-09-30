package commands

import (
	"dyc/internal/initialize"
	"dyc/internal/logger"
	"dyc/internal/module/article"
	"fmt"
	"github.com/urfave/cli"
	"strings"
	"sync"
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
	},
	Action: deployAction,
}

func deployAction(c *cli.Context) (err error) {
	initialize.Loading(c.String("conf"))
	dir := "/Users/liuning/Documents/github/book"
	conf, err := article.LoadDir(dir)
	if err != nil {
		logger.Fatalf("加载配置文件: %s", err)
	}
	// 清理一下文章
	if err = article.Purge(conf.Key); err != nil {
		logger.Fatalf("清理文章文件: %s", err)
	}
	// 初始化mapping
	if err = article.Initialize(); err != nil {
		logger.Fatalf("初始化es mapping: %s", err)
	}
	// 公众号二维码上传
	if err = conf.Qrcode(dir); err != nil {
		logger.Fatalf(": %s", err)
	}
	wg := sync.WaitGroup{}
	for t, i := range conf.Topics {
		for k, f := range i.Articles {
			wg.Add(1)
			logger.Debugf("analyze file: %s", f)
			go func(position int, title string, topic article.Topic, w *sync.WaitGroup) {
				defer w.Done()
				// 文件路径
				filename := fmt.Sprintf("/%s/%s/%s", strings.Trim(dir, "/"), strings.Trim(topic.Dir, "/"), strings.Trim(topic.Articles[position], "/"))
				// 图片资源路径
				assert := fmt.Sprintf("/%s/%s", strings.Trim(dir, "/"), strings.Trim(topic.Assert, "/"))
				a, err := article.NewArticle(filename)
				if err != nil {
					logger.Errorf("文章初始化失败: %s", err)
					return
				}
				// 数据完善
				a.Complete(conf, title, position)
				// 上传图片
				if err = a.UploadImage(assert, topic.Assert); err != nil {
					logger.Errorf("图片上传失败: %s", err)
					return
				}
				if err := a.Storage(); err != nil {
					logger.Errorf("elasticsearch 存储失败: %s", err)
					return
				}
			}(k, t, i, &wg)
		}
	}
	wg.Wait()
	return nil
}
