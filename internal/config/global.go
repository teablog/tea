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
