package deploy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/article"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

func Run(dir string) {
	logger.SetLevel("info")
	if ok, err := helper.Git.HasCommit(dir); err != nil {
		logger.Fatalf("check git commit dir: %s err: %s", dir, err.Error())
	} else if !ok {
		logger.Fatalf("please git commit the changes, dir: %s", dir)
	}
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
	all, err := article.Art.All([]string{"md5", "id", "title", "pv"})
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
	// error
	errs := make([]error, 0)
	chanErr := make(chan error, 2)
	defer close(chanErr)
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
	go func() {
		for {
			select {
			case err, ok := <-chanErr:
				if !ok {
					return
				}
				errs = append(errs, err)
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
					chanErr <- errors.Wrapf(err, "file load")
					return
				}
				// 数据完善
				if err := a.Complete(conf, topicTitle, file); err != nil {
					chanErr <- errors.Wrapf(err, "《%s》 complete %s", a.Title, err.Error())
					return
				}
				chanA <- a.ID
				// 文件没有变动
				if _, ok := md5Cache[a.Md5]; ok {
					return
				}
				// 上传图片
				if err = a.UploadImage(dir, a.Topic); err != nil {
					chanErr <- errors.Wrapf(err, "upload image: %s", err)
					return
				}
				// 新文章
				if _, ok := idCache[a.ID]; !ok {
					if err := a.Create(); err != nil {
						chanErr <- errors.Wrapf(err, "elasticsearch save err: %s", err)
						return
					}
				} else {
					if err := a.Update(); err != nil {
						chanErr <- errors.Wrapf(err, "elasticsearch save err: %s", err)
						return
					}
				}
			}(topicTitle, file)
		}
	}
	wg.Wait()
	close(chanA)
	// 等chanA goroutine退出以后在执行, 防止 concurrent map read and map write
	time.Sleep(10 * time.Millisecond)
	logger.Infof("----------------- errors --------------------")
	for _, v := range errs {
		logger.Errorf("%s", v.Error())
	}
	if len(errs) > 0 {
		logger.Error("please solve the errors")
		return
	}
	logger.Infof("----------------- delete article -----------")
	toDel := make([]string, 0)
	for _, v := range all {
		if _, ok := newCache[v.Id]; !ok {
			logger.Infof("delete: 《%s》pv %d id %s", v.Title, v.Pv, v.Id)
			toDel = append(toDel, v.Id)
		}
	}
	if err := article.Art.DeleteByIds(toDel); err != nil {
		logger.Errorf("delete err: %s", err.Error())
	}

	logger.Infof("--------- start convert webp ---------------")
	// 生成webp图片
	if err := helper.Image.Convert(path.Join(config.GetKey("path::storage_dir").String(), "images/blog", conf.Key)); err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("---------------- ping seo -----------------------")
	if err := regenSitemap(); err != nil {
		logger.Error(err)
	} else {
		if err := pingGoogleSitemap(); err != nil {
			logger.Error(err)
		}
	}
	if err := pingBaidu(); err != nil {
		logger.Error(err)
	}
}

func regenSitemap() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	clt := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	resp, err := clt.Get(config.Global.Host() + "/api/seo/sitemap")
	if err != nil {
		return errors.Wrapf(err, "regenerate sitemap err")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("regenerate sitemap failed!")
	}
	type r struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	re := new(r)
	if err := json.NewDecoder(resp.Body).Decode(re); err != nil {
		return errors.Wrap(err, "json decode response body err")
	}
	if re.Code != 0 {
		return errors.Errorf("regenerate sitemap failed")
	}
	return nil
}

// getSitemapUrl 获取sitemap地址
func getSitemapUrl() string {
	return config.Global.Host() + "/sitemap.xml"
}

// pingGoogleSitemap 向google提交站点地图 https://developers.google.com/search/docs/guides/submit-URLs?hl=zh-Hans
func pingGoogleSitemap() error {
	proxy, _ := url.Parse(config.Proxy.Http())
	clt := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	sitemapUrl := url.QueryEscape(getSitemapUrl())
	u := fmt.Sprintf("http://www.google.com/ping?sitemap=%s", sitemapUrl)
	resp, err := clt.Get(u)
	if err != nil {
		return errors.Wrapf(err, "ping google sitemap err")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("ping google sitemap failed")
	}
	logger.Info("ping google sitemap success")
	return nil
}

func pingGoogleSitemapProxySocks5() error {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:7890", nil, proxy.Direct)
	if err != nil {
		return errors.Wrapf(err, "socks5 direct err")
	}
	clt := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial:            dialer.Dial,
		},
	}
	sitemapUrl := url.QueryEscape(getSitemapUrl())
	u := fmt.Sprintf("http://www.google.com/ping?sitemap=%s", sitemapUrl)
	resp, err := clt.Get(u)
	if err != nil {
		return errors.Wrapf(err, "ping google sitemap err")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("ping google sitemap failed")
	}
	logger.Info("ping google sitemap success")
	return nil
}

func pingBaidu() error {
	clt := http.Client{
		Timeout: 5 * time.Second,
	}
	articles, err := article.Art.Today([]string{"id", "last_edit_time"})
	if err != nil {
		return err
	}
	if len(articles) < 0 {
		return errors.New("no articles")
	}
	buf := bytes.NewBuffer(nil)
	for _, v := range articles {
		u := fmt.Sprintf("%s/article/%s", config.Global.Host(), v.Id)
		buf.WriteString(u)
		buf.WriteString("\n")
	}
	u := "http://data.zz.baidu.com/urls?site=https://www.douyacun.com&token=mLTCWuzMZLOHBTYC"
	resp, err := clt.Post(u, "text/plain", buf)
	if err != nil {
		return errors.Wrapf(err, "ping baidu urls err")
	}
	defer resp.Body.Close()
	type r struct {
		Remain  int `json:"remain"`
		Success int `json:"success"`
	}
	res := new(r)
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return errors.Wrapf(err, "ping baidu urls err")
	}
	if res.Success > 0 {
		logger.Infof("ping baidu url success %d remain %d", res.Success, res.Remain)
		return nil
	} else {
		return errors.New("ping baidu urls success 0")
	}
}
