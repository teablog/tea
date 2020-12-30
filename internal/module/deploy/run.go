package deploy

import (
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/article"
	"path"
	"strings"
	"sync"
)

func Run(dir string) {
	//logger.SetLevel("info")
	conf, err := LoadConfig(dir)
	if err != nil {
		logger.Fatalf("load config file err: %s", err)
	}
	all, err := article.Post.All()
	if err != nil {
		logger.Fatalf("all articles md5 load err: %s", err.Error())
	}
	// md5
	md5Cache := all.MapMd5()
	idCache := all.MapId()

	if err = Indices.Article.Init(consts.IndicesArticleCost); err != nil {
		logger.Fatalf("article index init err: %s", err)
	}
	// 公众号二维码上传
	if err = conf.UploadQrcode(conf.Root); err != nil {
		logger.Fatalf("upload qr code err: %s", err)
	}
	wg := sync.WaitGroup{}
	// 加锁，限制并发数量
	queue := make(chan struct{}, 10)
	defer close(queue)
	for topicTitle, articles := range conf.Topics {
		for _, file := range articles {
			wg.Add(1)
			queue <- struct{}{}
			go func(topicTitle, file string) {
				defer wg.Done()
				defer func() {
					<-queue
				}()
				logger.Infof("start load file: %s", file)
				// 文件路径
				filePath := path.Join(dir, strings.ToLower(topicTitle), file)
				a, err := NewArticle(filePath)
				if err != nil {
					logger.Errorf("file load err: %s", err)
					return
				}
				// 文件没有变动
				if _, ok := md5Cache[a.Md5]; ok {
					return
				}
				// 数据完善
				a.Complete(conf, topicTitle, file)
				// 上传图片
				if err = a.UploadImage(dir, a.Topic); err != nil {
					logger.Errorf("upload image: %s", err)
					return
				}
				// 新文章
				if _, ok := idCache[a.ID]; !ok {
					if err := a.Create(); err != nil {
						logger.Errorf("elasticsearch save err: %s", err)
						return
					}
				} else {
					if err := a.Update(); err != nil {
						logger.Errorf("elasticsearch save err: %s", err)
						return
					}
				}
			}(topicTitle, file)
		}
	}
	wg.Wait()
	logger.Debugf("--------- start convert webp ---------------")
	// 生成webp图片
	if err := helper.Image.Convert(path.Join(config.GetKey("path::storage_dir").String(), "images/blog", conf.Key)); err != nil {
		logger.Error(err)
	}
}
