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
	"time"
)

func Run(dir string) {
	//logger.SetLevel("info")
	conf, err := LoadConfig(dir)
	if err != nil {
		logger.Fatalf("load config file err: %s", err)
	}
	if err = Indices.Article.Init(consts.IndicesArticleCost); err != nil {
		logger.Fatalf("article index init err: %s", err)
	}
	// 公众号二维码上传
	if err = conf.UploadQrcode(conf.Root); err != nil {
		logger.Fatalf("upload qr code err: %s", err)
	}
	// 历史文章
	all, err := article.Post.All()
	if err != nil {
		logger.Fatalf("all articles md5 load err: %s", err.Error())
	}
	// 历史数据
	md5Cache := all.MapMd5()
	idCache := all.MapId()
	newCache := make(map[string]struct{})
	wg := sync.WaitGroup{}
	// 锁：限制并发数量
	concurrentN := make(chan struct{}, 10)
	// chan: 统计当前文章
	chanA := make(chan string, 2)
	defer close(concurrentN)

	go func() {
		for {
			select {
			case id, ok := <-chanA:
				if !ok {
					return
				}
				newCache[id] = struct{}{}
			}
		}
	}()
	for topicTitle, articles := range conf.Topics {
		for _, file := range articles {
			wg.Add(1)
			concurrentN <- struct{}{}
			go func(topicTitle, file string) {
				defer wg.Done()
				defer func() {
					<-concurrentN
				}()
				logger.Infof("start load file: %s", file)
				// 文件路径
				filePath := path.Join(dir, strings.ToLower(topicTitle), file)
				a, err := NewArticle(filePath)
				if err != nil {
					logger.Errorf("file load err: %s", err)
					return
				}
				// 数据完善
				if err := a.Complete(conf, topicTitle, file); err != nil {
					logger.Errorf("《%s》 complete %s", a.Title, err.Error())
				}
				chanA <- a.ID
				// 文件没有变动
				if _, ok := md5Cache[a.Md5]; ok {
					return
				}
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
	close(chanA)
	// 等chanA goroutine推出以后在执行, 防止 concurrent map read and map write
	time.Sleep(10 * time.Millisecond)
	logger.Infof("----------------- delete article -----------")
	toDel := make([]string, 0)
	for _, v := range all {
		if _, ok := newCache[v.Id]; !ok {
			logger.Infof("Delete: 《%s》pv %d id %s", v.Title, v.Pv, v.Id)
			toDel = append(toDel, v.Id)
		}
	}
	if err := article.Post.DeleteByIds(toDel); err != nil {
		logger.Errorf("[delete] err: %s", err.Error())
	}

	logger.Infof("--------- start convert webp ---------------")
	// 生成webp图片
	if err := helper.Image.Convert(path.Join(config.GetKey("path::storage_dir").String(), "images/blog", conf.Key)); err != nil {
		logger.Error(err)
	}
}
