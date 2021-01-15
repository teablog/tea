package logstash

import (
	"fmt"
	"github.com/teablog/tea/internal/config"
	"testing"
)

func TestMatch(t *testing.T) {
	config.Init("configs/debug.ini")
	s := "Mozilla/5.0 (compatible; Baiduspider/2.0; Baiduspider+http://www.baidu.com/search/spider.html)"
	fmt.Println(match(s))
}
