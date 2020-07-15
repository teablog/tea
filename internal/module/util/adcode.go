package util

import (
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var AdCoder *AdCode

type AdCode struct {
	Name     string `json:"name"`
	Adcode   string `json:"adcode"`
	CityCode string `json:"city_code"`
}

func (*AdCode) FindByName(ctx *gin.Context, name string) (*[]AdCode, error) {
	body := fmt.Sprintf(`{
  "query": {
    "match": {
      "name": "%s"
    }
  },
  "size": 5
}`, name)
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesAdCodeConst),
		db.ES.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, errors.New(string(bodyRaw))
	}
	var r db.ESListResponse
	if err = json.Unmarshal(bodyRaw, &r); err != nil {
		return nil, err
	}
	if r.Hits.Total.Value == 0 {
		return nil, nil
	}
	var hits []db.ESItemResponse
	if err = json.Unmarshal(r.Hits.Hits, &hits); err != nil {
		return nil, err
	}
	var list []AdCode
	for _, v := range hits {
		var source AdCode
		if err = json.Unmarshal(v.Source, &source); err == nil {
			list = append(list, source)
		}
	}
	return &list, nil
}

func (a *AdCode) FindCity(ctx *gin.Context, name string) (*AdCode, error) {
	if regions, err := a.FindByName(ctx, name); err != nil {
		return nil, err
	} else {
		for _, v := range *regions {
			if v.IsCity(v.Adcode) && strings.HasPrefix(v.Name, name) {
				return &v, nil
			}
		}
		return nil, errors.New("city not found")
	}
}

func (*AdCode) FindByCode(ctx *gin.Context, code string) (*AdCode, error) {
	query := fmt.Sprintf(`{
  "query": {
    "term": {
      "adcode": "%s"
    }
  }
}`, code)
	if res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesAdCodeConst),
		db.ES.Search.WithBody(strings.NewReader(query)),
	); err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		raw, _ := ioutil.ReadAll(res.Body)
		if res.IsError() {
			return nil, errors.New(string(raw))
		}
		var r db.ESListResponse
		if err = json.Unmarshal(raw, &r); err != nil {
			return nil, err
		}
		if r.Hits.Total.Value == 0 {
			return nil, nil
		}
		var hits []db.ESItemResponse
		if err = json.Unmarshal(r.Hits.Hits, &hits); err != nil {
			return nil, err
		}
		var source AdCode
		if err = json.Unmarshal(hits[0].Source, &source); err != nil {
			return nil, err
		}
		return &source, nil
	}
}

func (a *AdCode) BelongProvince(ctx *gin.Context, code string) (*AdCode, error) {
	return a.FindByCode(ctx, code[:2]+"0000")
}

func (a *AdCode) BelongCity(ctx *gin.Context, code string) (*AdCode, error) {
	return a.FindByCode(ctx, code[:4]+"00")
}

func (*AdCode) IsProvince(code string) bool {
	return code[2:] == "0000"
}

func (*AdCode) IsCity(code string) bool {
	return code[2:4] != "00" && code[4:] == "00"
}

func (*AdCode) IsDistrict(code string) bool {
	return code[2:4] != "00" && code[4:] != "00"
}

func (*AdCode) CanFindCity(code string) bool {
	return code[2:4] != "00"
}

func (a *AdCode) Component(ctx *gin.Context, code string) (res map[string]*AdCode, err error) {
	res = make(map[string]*AdCode)
	if res["province"], err = a.BelongProvince(ctx, code); err != nil {
		return
	}
	if a.CanFindCity(code) {
		if res["city"], err = a.BelongCity(ctx, code); err != nil {
			return
		} else {
			if strings.Contains(res["city"].Name, "市辖区") {
				res["city"].Name = res["city"].Name[:len(res["city"].Name)-9]
			}
		}
	}
	if a.IsDistrict(code) {
		if res["district"], err = a.FindByCode(ctx, code); err != nil {
			return
		}
	}
	return
}
