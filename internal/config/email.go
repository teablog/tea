package config

var Email *_email

type _email struct{}

func (*_email) Enable() bool {
	ok, _ := Config.Section("email").Key("enable").Bool()
	return ok
}

func (*_email) Username() string {
	return Config.Section("email").Key("username").String()
}

func (*_email) Password() string {
	return Config.Section("email").Key("password").String()
}

func (*_email) Host() string {
	return Config.Section("email").Key("host").String()
}

func (*_email) Port() int {
	port, _ := Config.Section("email").Key("port").Int()
	return port
}
