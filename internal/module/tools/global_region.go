package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	GlobalRegionLevelContinent = iota + 1
	GlobalRegionLevelCountry
	GlobalRegionLevelProvince
	GlobalRegionLevelCity
	GlobalRegionLevelDistinct
)

var CityNotFound = errors.New("city not found")
var CountryNotFound = errors.New("country not found")

var GlobalRegion *_globalRegion

type _globalRegion GlobalRegionComponent

type GlobalRegionComponent struct {
	Id         int    `json:"id"`
	Pid        int    `json:"pid"`
	Path       string `json:"path"`
	Level      int    `json:"level"`
	Name       string `json:"name"`
	NameEN     string `json:"name_en"`
	NamePinyin string `json:"name_pinyin"`
	Code       string `json:"code"`
}

func (g *_globalRegion) FindByNameEN(ctx context.Context, nameEN string) ([]*GlobalRegionComponent, error) {
	query := fmt.Sprintf(`
{
  "query": {
    "term": {
      "name_en": {
        "value": "%s"
      }
    }
  }
}`, nameEN)
	resp, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesGlobalRegion),
		db.ES.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "es search failed")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.IsError() {
		return nil, errors.Errorf("es search error: %s", body)
	}
	r := &db.ESListResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "es list response error:")
	}
	hits := make([]db.ESItemResponse, 0)
	if err := json.Unmarshal(r.Hits.Hits, &hits); err != nil {
		return nil, errors.Wrap(err, "es item response error:")
	}
	rows := make([]*GlobalRegionComponent, 0, len(hits))
	if len(hits) == 0 {
		return rows, nil
	}
	for _, v := range hits {
		item := &GlobalRegionComponent{}
		if err := json.Unmarshal(v.Source, item); err != nil {
			rows = append(rows, item)
		}
	}
	return rows, nil
}

func (g *_globalRegion) GetCityByNameEN(ctx context.Context, nameEN string) (*GlobalRegionComponent, error) {
	rows, err := g.FindByNameEN(ctx, nameEN)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, CityNotFound
	}
	city := rows[0]
	if city.Level != GlobalRegionLevelCity {
		return nil, CityNotFound
	}
	return city, nil
}

func (g *_globalRegion) GetCountryByNameEN(ctx context.Context, nameEN string) (*GlobalRegionComponent, error) {
	rows, err := g.FindByNameEN(ctx, nameEN)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, CountryNotFound
	}
	country := rows[0]
	if country.Level != GlobalRegionLevelCountry {
		return nil, CountryNotFound
	}
	return country, nil
}

func (g *_globalRegion) FindByID(ctx context.Context, id string) (*GlobalRegionComponent, error) {
	resp, err := db.ES.Get(consts.IndicesGlobalRegion, id)
	if err != nil {
		return nil, errors.Wrap(err, "global region find error")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.IsError() {
		return nil, errors.Errorf("[%s] find failed, error: %s", id, body)
	}
	r := &db.ESItemResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "es item response:")
	}
	region := &GlobalRegionComponent{}
	if err := json.Unmarshal(r.Source, region); err != nil {
		return nil, errors.Wrap(err, "global region component:")
	}
	return region, nil
}

func (g *_globalRegion) Component(ctx context.Context, cityId string) (map[string]*AdCodeComponent, error) {
	if city, err := g.FindByID(ctx, cityId); err != nil {
		return nil, err
	} else {
		if country, err := g.FindByID(ctx, strconv.Itoa(city.Pid)); err != nil {
			return nil, err
		} else {
			return map[string]*AdCodeComponent{
				"country": {
					Name:   country.NameEN,
					Adcode: "",
					Code:   country.Code,
				},
				"city": {
					Name:   city.Name,
					Adcode: "",
					Code:   city.Code,
				},
			}, nil
		}
	}
}
