package deploy

import (
	"dyc/internal/helper"
	"dyc/internal/logger"
	"sync"
)

func Run(dir string) {
	conf, err := LoadConfig(dir)
	if err != nil {
		logger.Fatalf("加载配置文件: %s", err)
	}
	// 清理一下文章
	if err = Topic.Init(); err != nil {
		logger.Fatalf("初始化: %s", err)
	}
	// 初始化mapping
	//if err = Topic.Purge(conf.Key); err != nil {
	//	logger.Fatalf("初始化es mapping: %s", err)
	//}
	// 公众号二维码上传
	if err = conf.UploadQrcode(conf.Root); err != nil {
		logger.Fatalf(": %s", err)
	}
	wg := sync.WaitGroup{}
	for topicTitle, articles := range conf.Topics {
		for _, file := range articles {
			wg.Add(1)
			logger.Debugf("analyze file: %s", file)
			//
			go func(topicTitle, file string, w *sync.WaitGroup) {
				defer w.Done()
				// 文件路径
				filePath := helper.Path.Join(dir, topicTitle, file)
				a, err := NewArticle(filePath)
				if err != nil {
					logger.Errorf("文章初始化失败: %s", err)
					return
				}
				// 数据完善
				a.Complete(conf, topicTitle, file)
				// 上传图片
				if err = a.UploadImage(helper.Path.Join(dir, a.Topic)); err != nil {
					logger.Errorf("图片上传失败: %s", err)
					return
				}
				if err := a.Storage(); err != nil {
					logger.Errorf("elasticsearch 存储失败: %s", err)
					return
				}
			}(topicTitle, file, &wg)
		}
	}
	wg.Wait()
}
