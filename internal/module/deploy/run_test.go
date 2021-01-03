package deploy

import (
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/logger"
	"os"
	"testing"
)

func TestPingGoogleSitemap(t *testing.T) {
	config.Init("configs/prod.ini")
	logger.NewDefaultLogger(os.Stdout)
	t.Log(config.Proxy.Http())
	if err := pingGoogleSitemapProxySocks5(); err != nil {
		t.Error(err)
	}
}
