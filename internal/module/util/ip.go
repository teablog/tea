package util

import (
	"github.com/teablog/tea/internal/config"
	"dyc/internal/module/util/ip2location"
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
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
)

func Init() {
	ipipdb, _ = ipdb.NewCity(config.GetKey("ip::ipip_file").String())
	ip2locationdb, _ = ip2location.OpenDB(config.GetKey("ip::ip2location_file").String())
}

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
