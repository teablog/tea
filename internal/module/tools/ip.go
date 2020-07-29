package tools

import (
	"context"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/module/tools/ip2location"
)

var (
	LocationError = errors.New("locate failed")
	ipipdb        *ipdb.City
	ip2locationdb *ip2location.DB
)

const (
	Language = "zh-CN"

	ReferIp2location = "ip2location"
	ReferIpIp        = "ipip"

	ChinaEN = "China"
)

type Region struct {
	Country  string
	Province string
	City     string
	District string
	Refer    string
}

func Init() {
	var err error
	if ipipdb, err = ipdb.NewCity(config.GetKey("ip::ipip_file").String()); err != nil {
		panic(errors.Wrap(err, "ipip db init error"))
	}
	if ip2locationdb, err = ip2location.OpenDB(config.GetKey("ip::ip2location_file").String()); err != nil {
		panic(errors.Wrap(err, "ip2location db init error"))
	}
}

func ipip2(ip string) (*Region, error) {
	info, err := ipipdb.FindMap(ip, "CN")
	if err != nil {
		return nil, err
	}
	return &Region{
		Province: info["country_name"],
		City:     info["city_name"],
		District: "",
		Refer:    ReferIpIp,
	}, nil
}

func ip2location2(ip string) (*Region, error) {
	res, err := ip2locationdb.Get_all(ip)
	if err != nil {
		return nil, errors.Wrap(err, "ip2location get all error")
	}
	if res.City == "-" || res.Country_long == "-" {
		return nil, CityNotFound
	}
	return &Region{
		Country:  res.Country_long,
		Province: res.Region,
		City:     res.City,
		District: "",
		Refer:    ReferIp2location,
	}, nil
}

func LocationByIp(ctx context.Context, ip string) (map[string]*AdCodeComponent, error) {
	if res, err := ipip2(ip); err == nil && res.City != "" {
		if city, err := AdCode.FindCity(ctx, res.City); err == nil {
			return AdCode.Component(ctx, city.Adcode)
		}
	}
	if res, err := ip2location2(ip); err == nil {
		if res.Country == ChinaEN {
			if region, err := GlobalRegion.GetChinaCityByNameEN(ctx, res.City); err == nil {
				if city, err := AdCode.FindCity(ctx, region.Name); err == nil {
					return AdCode.Component(ctx, city.Adcode)
				}
			}
		} else {
			//if region, err := GlobalRegion.GetCityByNameEN(ctx, res.City); err == nil {
			//	return GlobalRegion.Component(ctx, strconv.Itoa(region.Id))
			//}
			return map[string]*AdCodeComponent{
				"country": {
					Name: res.Country,
				},
				"city": {
					Name: res.City,
				},
			}, nil
		}
	}
	return nil, LocationError
}
