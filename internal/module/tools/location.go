package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var Location location

type location struct{}
type amapString string

func (a *amapString) UnmarshalJSON(b []byte) error {
	if string(b) == "[]" {
		*a = ""
	} else {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*a = amapString(s)
	}
	return nil
}

type addressComponent struct {
	Province     amapString `json:"province"`
	City         amapString `json:"city"`
	CityCode     amapString `json:"citycode"`
	District     amapString `json:"district"`
	AdCode       amapString `json:"adcode"`
	Towncode     amapString `json:"towncode"`
	Neighborhood struct {
		Name amapString `json:"name"`
		Type amapString `json:"type"`
	} `json:"neighborhood"`
	Building struct {
		Name amapString `json:"name"`
		Type amapString `json:"type"`
	} `json:"building"`
	StreetNumber struct {
		Street    amapString `json:"street"`
		Number    amapString `json:"number"`
		Location  amapString `json:"location"`
		Direction amapString `json:"direction"`
		Distance  amapString `json:"distance"`
	} `json:"streetNumber"`
}

type regeo struct {
	Status    json.Number `json:"status"`
	Info      string      `json:"info"`
	Regeocode struct {
		FormattedAddress string           `json:"formatted_address"`
		AddressComponent addressComponent `json:"addressComponent"`
	} `json:"regeocode"`
}

func (*location) FindByGeoCode(ctx context.Context, latitude, longitude string) (*addressComponent, error) {
	url := "https://restapi.amap.com/v3/geocode/regeo?key=959aef6c8cf282cf66640c5ad83d8298&location=%s,%s&extensions=base"
	resp, err := http.Get(fmt.Sprintf(url, latitude, longitude))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r regeo
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	address := &r.Regeocode.AddressComponent

	return address, err
}
