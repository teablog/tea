package config

import "fmt"

var Global *_global

type _global struct{}

// Host 获取完整域名
func (g *_global) Host() string {
	return fmt.Sprintf("%s://%s", g.Protocol(), g.Domain())
}

// Domain 域名
func (*_global) Domain() string {
	return Config.Section("global").Key("domain").String()
}

// Protocol 协议
func (*_global) Protocol() string {
	return Config.Section("global").Key("protocol").String()
}

func (*_global) CdnHost() string {
	return Config.Section("global").Key("cdn_host").String()
}

func (*_global) Token() string {
	return Config.Section("global").Key("token").String()
}

func (*_global) CookieMaxAge() int {
	t, _ := Config.Section("global").Key("cookie_max_age").Int()
	if t == 0 {
		return 31536000
	}
	return t
}
