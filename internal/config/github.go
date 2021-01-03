package config

var Github *_github

type _github struct{}

func (*_github) ClientId() string {
	return Config.Section("github").Key("client_id").String()
}

func (*_github) ClientSecret() string {
	return Config.Section("github").Key("client_secret").String()
}
