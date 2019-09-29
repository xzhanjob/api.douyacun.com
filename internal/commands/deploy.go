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
		logger.Fatalf("load dir failed: %s", err)
	}
	wg := sync.WaitGroup{}
	for t, i := range conf.Topics {
		for k, f := range i.Articles {
			wg.Add(1)
			logger.Debugf("analyze file: %s", f)
			go func(position int, title string, topic article.Topic, w *sync.WaitGroup) {
				defer w.Done()
				// 文件路径
				filename := fmt.Sprintf("/%s/%s/%s", strings.Trim(dir, "/"), strings.Trim(topic.Dir, "/"), strings.Trim(f, "/"))
				// 图片资源路径
				assert := fmt.Sprintf("/%s/%s", strings.Trim(dir, "/"), strings.Trim(topic.Assert, "/"))
				a, err := article.NewArticle(filename)
				if err != nil {
					logger.Errorf("new article failed: %s", err)
					return
				}
				// 数据完善
				a.Complete(conf, title, position)
				// 上传图片
				if err = a.UploadImage(assert, topic.Assert); err != nil {
					logger.Errorf("upload image failed: %s", err)
					return
				}
				if err := a.Storage(); err != nil {
					logger.Errorf("elasticsearch storage failed: %s", err)
					return
				}
			}(k, t, i, &wg)
		}
	}
	wg.Wait()
	return nil
}
