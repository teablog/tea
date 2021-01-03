package config

var Proxy *_proxy

type _proxy struct{}

func (p *_proxy) Http() string {
	return Config.Section("proxy").Key("http").String()
}

func (*_proxy) Socks5() string {
	return Config.Section("proxy").Key("socks5").String()
}
