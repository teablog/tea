package outside

import (
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/mail"
	"os"
	"testing"
	"time"
)

func TestSpider(t *testing.T) {
	config.Init("configs/debug.ini")
	logger.NewLogger(os.Stdout, "debug")
	mail.Init()
	db.NewElasticsearch(config.ES.Address(), config.ES.User(), config.ES.Password())
	o := new(Outside)
	o.Url = "http://douyacun.io"
	o.Email = "douyacun@163.com"
	o.Title = "Douyacun"
	o.Spider()
	time.Sleep(2 * time.Second)
}
