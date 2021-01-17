package config

var Ad *_ad

type _ad struct{}

// ad sense: 首页feeds流位置
func (*_ad) AdSenseFeedsPos() int {
	 i, _ := Config.Section("ad").Key("ad_sense_feeds_pos").Int()
	 return i
}
