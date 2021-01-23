package seo

import (
	"errors"
	"fmt"
	"github.com/douyacun/gositemap"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/module/article"
	"math"
	"path"
	"time"
)

var Sitemap sitemap

type sitemap struct{}

// sitemap.xml 导出
// 1. 从google导出已经收录的文章, 权重从3开始按时间降低
// 2. 按时间开始降权
// 3. sitemap.xml 收录图片
func (s *sitemap) Generate(ctx *gin.Context) error {
	if err := article.Art.Flush(); err != nil {
		return err
	}
	articles, err := article.Art.All([]string{"id", "last_edit_time"})
	if err != nil {
		return err
	}
	if len(articles) < 0 {
		return errors.New("no articles")
	}
	st := gositemap.NewSiteMap()
	st.SetPretty(true)
	st.SetCompress(false)
	st.SetDefaultHost(config.Global.Host())
	st.SetPublicPath(path.Join(config.Path.StorageDir(), "public"))
	host := config.Global.Host() + "/article/%s"

	url := gositemap.NewUrl()
	url.SetLoc(config.Global.Host())
	url.SetChangefreq(gositemap.Daily)
	url.SetPriority(1)
	st.AppendUrl(url)

	for _, v := range articles {
		url := gositemap.NewUrl()
		url.SetLoc(fmt.Sprintf(host, v.Id))
		url.SetLastmod(v.LastEditTime)
		url.SetPriority(s.priority(v.LastEditTime))
		url.SetChangefreq(s.freq(v.LastEditTime))
		st.AppendUrl(url)
	}
	_, err = st.Storage()
	if err != nil {
		return err
	}
	return nil
}

// 计算权重
// 30天权重减1
// 不足30天: (t % 30)/30
func (s *sitemap) priority(t time.Time) float64 {
	p := float64(1)   // 最大权重
	mp := 0.1         // 月权重
	min := float64(0) // 最小权重
	now := time.Now()
	d := now.Sub(t)
	h := d.Hours()
	m := math.Ceil(h / 720)
	p = p - mp*m
	if p < 0 {
		p = min
	}
	dp := float64(int64(h)%720)/720*mp
	p = p + dp
	return p
}

func (s *sitemap) freq(t time.Time) gositemap.ChangeFreq {
	h := time.Now().Sub(t).Hours()
	if h <= 168 {
		return gositemap.Daily
	} else if h > 168 && h < 720 {
		return gositemap.Weekly
	} else {
		return gositemap.Monthly
	}
}
