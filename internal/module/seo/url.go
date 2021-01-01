package seo

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/module/article"
	"os"
	"path"
)

var Url *_url

type _url struct{}

// Generate 百度seo生成url.txt https://ziyuan.baidu.com/college/articleinfo?id=3159
func (*_url) Generate(ctx *gin.Context) error {
	articles := article.Search.All([]string{"id", "last_edit_time"})
	if len(articles) < 0 {
		return errors.New("no articles")
	}
	buf := bytes.NewBuffer(nil)
	for _, v := range articles {
		u := fmt.Sprintf("%s/article/%s", config.Global.Host(), v.Id)
		buf.WriteString(u)
		buf.WriteString("\n")
	}
	p := config.Path.StorageDir()
	filepath := path.Join(p, "seo", "url.txt")
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "baidu url.txt open %s err", filepath)
	}
	if _, err = f.Write(buf.Bytes()); err != nil {
		return errors.Wrapf(err, "baidu url.txt write err")
	}
	return nil
}

