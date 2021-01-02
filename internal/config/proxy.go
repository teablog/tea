package config

import "fmt"

var Proxy *_proxy

type _proxy struct{}

func (p *_proxy) GetLocalEndpoint() string {
	return fmt.Sprintf("http://127.0.0.1:%s", p.GetEndpoint())
}

func (*_proxy) GetEndpoint() string {
	return Config.Section("proxy").Key("endpoint").String()
}
