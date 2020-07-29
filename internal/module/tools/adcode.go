package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"io/ioutil"
	"strings"
)

var AdCode *AdCodeComponent

type AdCodeComponent struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
	Code   string `json:"code"`
}

func (*AdCodeComponent) NewDefault() map[string]AdCodeComponent {
	return map[string]AdCodeComponent{
		"country": {
			Name:   "中国",
			Adcode: "",
			Code:   "CN",
		},
		"city": {
			Name:   "北京市",
			Adcode: "110100",
			Code:   "",
		},
		"province": {
			Name:   "北京市",
			Adcode: "110000",
			Code:   "",
		},
	}
}

func (*AdCodeComponent) FindByName(ctx context.Context, name string) (*[]AdCodeComponent, error) {
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
	var list []AdCodeComponent
	for _, v := range hits {
		var source AdCodeComponent
		if err = json.Unmarshal(v.Source, &source); err == nil {
			list = append(list, source)
		}
	}
	return &list, nil
}

func (*AdCodeComponent) FindByNamePingyin(ctx context.Context, pingyin string) (*[]AdCodeComponent, error) {
	body := fmt.Sprintf(`{
  "query": {
    "prefix": {
      "name.pinyin": "%s"
    }
  }
}`, strings.ToLower(pingyin))
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesAdCodeConst),
		db.ES.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("response body read failed")
	}
	if res.IsError() {
		return nil, errors.New(string(bodyRaw))
	}
	return nil, nil
}

func (a *AdCodeComponent) FindCity(ctx context.Context, name string) (*AdCodeComponent, error) {
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

func (*AdCodeComponent) FindByCode(ctx context.Context, code string) (*AdCodeComponent, error) {
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
		var source AdCodeComponent
		if err = json.Unmarshal(hits[0].Source, &source); err != nil {
			return nil, err
		}
		return &source, nil
	}
}

func (a *AdCodeComponent) BelongProvince(ctx context.Context, code string) (*AdCodeComponent, error) {
	return a.FindByCode(ctx, code[:2]+"0000")
}

func (a *AdCodeComponent) BelongCity(ctx context.Context, code string) (*AdCodeComponent, error) {
	return a.FindByCode(ctx, code[:4]+"00")
}

func (*AdCodeComponent) IsProvince(code string) bool {
	return code[2:] == "0000"
}

func (*AdCodeComponent) IsCity(code string) bool {
	return code[2:4] != "00" && code[4:] == "00"
}

func (*AdCodeComponent) IsDistrict(code string) bool {
	return code[2:4] != "00" && code[4:] != "00"
}

func (*AdCodeComponent) CanFindCity(code string) bool {
	return code[2:4] != "00"
}

func (a *AdCodeComponent) Component(ctx context.Context, code string) (res map[string]*AdCodeComponent, err error) {
	res = make(map[string]*AdCodeComponent)
	res["country"] = &AdCodeComponent{
		Name:   "中国",
		Adcode: "",
		Code:   "CN",
	}
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
