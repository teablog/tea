package deploy

import (
	"github.com/teablog/tea/internal/config"
	"testing"
)

func TestPingGoogleSitemap(t *testing.T) {
	config.Init("configs/prod.ini")
	t.Log(config.Proxy.GetLocalEndpoint())
	if err := pingGoogleSitemap(); err != nil {
		t.Error(err)
	}
}
