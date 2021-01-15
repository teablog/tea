package logstash

import (
	"fmt"
	"github.com/teablog/tea/internal/config"
	"testing"
)

func TestMatch(t *testing.T) {
	config.Init("configs/debug.ini")
	s := "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)"
	fmt.Println(match(s))
}
