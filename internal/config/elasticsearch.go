package config

var ES *_elasticsearch

type _elasticsearch struct{}

func (*_elasticsearch) Address() []string {
	return Config.Section("elasticsearch").Key("address").Strings(",")
}

func (*_elasticsearch) User() string {
	return Config.Section("elasticsearch").Key("user").String()
}

func (*_elasticsearch) Password() string {
	return Config.Section("elasticsearch").Key("password").String()
}

func (*_elasticsearch) FriendsLinkId() string {
	return Config.Section("elasticsearch").Key("friends_link_id").String()
}
