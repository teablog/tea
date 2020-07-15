package article

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"github.com/pkg/errors"
)

var Label _labels

type _labels struct {
	Label string `yaml:"label" json:"label"`
}

func (*_labels) List(size int) (data []string, err error) {
	var (
		buf bytes.Buffer
	)
	query := map[string]interface{}{
		"_source": []string{"label"},
		"size":    size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": map[string]interface{}{
					"script": map[string]interface{}{
						"script": map[string]interface{}{
							"source": "doc['label'].size() > 0",
						},
					},
				},
			},
		},
	}
	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		panic(errors.Wrap(err, "json encode错误"))
	}
	if res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(&buf),
	); err != nil {
		panic(errors.Wrap(err, "es search error"))
	} else {
		defer res.Body.Close()
		if res.IsError() {
			panic(errors.Wrap(err, "es 查询错误")) // 如果label为空直接panic，定位问题然后修复
		}
		var r db.ESListResponse
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			panic(errors.Wrap(err, "json decode错误"))
		}
		type _source struct {
			Label string `json:"label"`
		}
		var hits []db.ESItemResponse
		if err = json.Unmarshal(r.Hits.Hits, &hits); err != nil {
			return nil, err
		}
		for _, v := range hits {
			var source _source
			if err := json.Unmarshal(v.Source, &source); err != nil {
				panic(err)
			}
			data = append(data, source.Label)
		}
		return data, nil
	}
}
