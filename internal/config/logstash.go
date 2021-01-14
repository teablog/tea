package config

var Logstash *_logstash

type _logstash struct{}

func (*_logstash) Spider() string {
	return Config.Section("logstash").Key("spider").String()
}
