package seo

import (
	"errors"
	"fmt"
	"github.com/douyacun/gositemap"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/module/article"
	"path"
)

var Sitemap sitemap

type sitemap struct{}

func (s *sitemap) Generate(ctx *gin.Context) error {
	if err := article.Post.Flush(); err != nil {
		return err
	}
	articles, err := article.Post.All([]string{"id", "last_edit_time"})
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
	st.SetPublicPath(path.Join(config.Path.StorageDir(), "seo"))
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
		url.SetPriority(0.8)
		url.SetChangefreq(gositemap.Monthly)
		st.AppendUrl(url)
	}
	_, err = st.Storage()
	if err != nil {
		return err
	}
	return nil
}
