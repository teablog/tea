package util

import (
	"github.com/teablog/tea/internal/config"
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/oschwald/geoip2-golang"
	"net"
)

var (
	ipipdb  *ipdb.City
	geoipdb *geoip2.Reader
)

const (
	Language = "zh-CN"
)

func Init() {
	ipipdb, _ = ipdb.NewCity(config.GetKey("ip::ipip_file").String())
	geoipdb, _ = geoip2.Open(config.GetKey("ip::geo_file").String())
}

func ipip(ip string) (map[string]string, error) {
	info, err := ipipdb.FindMap(ip, "CN")
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	res["city"] = info["city_name"]
	res["country"] = info["country_name"]
	return res, nil
}

func geoip(ip string) (map[string]string, error) {
	record, err := geoipdb.City(net.ParseIP(ip))
	if err != nil {
		return nil, err
	}
	_ = record.City.Names
	res := make(map[string]string)
	res["city"] = record.City.Names[Language]
	res["continent"] = record.Continent.Names[Language]
	res["country"] = record.Country.Names[Language]
	res["latitude"] = fmt.Sprintf("%f", record.Location.Latitude)
	res["longitude"] = fmt.Sprintf("%f", record.Location.Latitude)
	return res, nil
}

func LocationByIp(ip string) (map[string]string, error) {
	res, err := ipip(ip)
	if err != nil || res["city"] == "" {
		res, err = geoip(ip)
	}
	return res, err
}

func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
