package config

var Sever *_server

type _server struct{}

func (*_server) CdnHost() string {
	return Config.Section("server").Key("cdn_host").String()
}
