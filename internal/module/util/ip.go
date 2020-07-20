package util

import (
	"dyc/internal/config"
	"dyc/internal/module/util/ip2location"
	"dyc/internal/module/util/qqwry"
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
	"net"
	"os"
)

var (
	ipipdb        *ipdb.City
	geoipdb       *geoip2.Reader
	qqwrydb       *os.File
	ip2locationdb *ip2location.DB
)

const (
	Language = "zh-CN"

	ReferIp2location = "ip2location"
	ReferIpIp        = "ipip"
	ReferQQwry       = "qqwry"
	ReferGeoIp       = "geoip"
)

func Init() {
	ipipdb, _ = ipdb.NewCity(config.GetKey("ip::ipip_file").String())
	geoipdb, _ = geoip2.Open(config.GetKey("ip::geo_file").String())
	qqwrydb, _ = qqwry.Getqqdata(config.GetKey("ip::qqwry_file").String())
	ip2locationdb, _ = ip2location.OpenDB(config.GetKey("ip::ip2location_file").String())
}

type parseIpFunc func(ip string) (map[string]string, error)

func ipip2(ip string) (map[string]string, error) {
	info, err := ipipdb.FindMap(ip, "CN")
	if err != nil {
		return nil, err
	}
	row := make(map[string]string)
	row["city"] = info["city_name"]
	row["country"] = info["country_name"]
	row["refer"] = ReferIpIp
	return row, nil
}

func geoip(ip string) (map[string]string, error) {
	record, err := geoipdb.City(net.ParseIP(ip))
	if err != nil {
		return nil, err
	}
	_ = record.City.Names
	row := make(map[string]string)
	row["city"] = record.City.Names[Language]
	row["continent"] = record.Continent.Names[Language]
	row["country"] = record.Country.Names[Language]
	row["latitude"] = fmt.Sprintf("%f", record.Location.Latitude)
	row["longitude"] = fmt.Sprintf("%f", record.Location.Latitude)
	row["refer"] = ReferGeoIp
	return row, nil
}

func qqwry2(ip string) (map[string]string, error) {
	var (
		row = make(map[string]string)
	)
	row["country"], row["city"] = qqwry.Getlocation(qqwrydb, ip)
	row["refer"] = ReferQQwry
	return row, nil
}

func ip2location2(ip string) (map[string]string, error) {
	var row = make(map[string]string)
	res, err := ip2locationdb.Get_all(ip)
	if err != nil {
		return nil, errors.Wrap(err, "ip2location get all error")
	}
	row["country"] = res.Country_long
	row["region"] = res.Region
	row["city"] = res.City
	row["refer"] = ReferIp2location
	return row, nil
}

func workflow(call parseIpFunc, ch chan map[string]string, ip string) {
	res, err := call(ip)
	if err != nil {
		return
	}
	ch <- res
}

func LocationByIp(ip string) (map[string]string, error) {
	res, err := ip2location2(ip)
	//if err != nil || res["city"] == "" {
	//	res, err = geoip(ip)
	//}

	//ch := make(chan map[string]string, 4)
	//
	//wg := &sync.WaitGroup{}
	//wg.Add(4)
	//go workflow(ipip2, ch, ip)
	//go workflow(geoip, ch, ip)
	//go workflow(ip2location2, ch, ip)
	//go workflow(qqwry2, ch, ip)
	//wg.Wait()
	//close(ch)
	//for v := range ch {
	//
	//}

	return res, err
}



func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
