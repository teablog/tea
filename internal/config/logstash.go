package config

import "strings"

var Logstash *_logstash

type _logstash struct{}

func (*_logstash) Spider() string {
	return Config.Section("logstash").Key("spider").String()
}

func (*_logstash) FilterUri() []string {
	uri := Config.Section("logstash").Key("filter_uri").String()
	return strings.Split(uri, ",")
}
