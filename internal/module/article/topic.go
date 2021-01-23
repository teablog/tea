package article

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
)

var Topics _topic

type _topic struct{}

func (*_topic) List(topic string, page int) (total int, data ASlice, err error) {
	var (
		buf bytes.Buffer
	)
	data = make(ASlice, 0)
	skip := (page - 1) * PageSize
	query := map[string]interface{}{
		"from": skip,
		"size": PageSize,
		"sort": map[string]interface{}{
			"last_edit_time": map[string]interface{}{
				"order": "desc",
			},
		},
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"topic.keyword": map[string]interface{}{
					"value": topic,
				},
			},
		},
		"_source": []string{"author", "title", "description", "topic", "id", "cover", "date", "last_edit_time"},
	}
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode错误"))
	}
	total, data, err = Art.Search(buf.String())
	adPos := config.Ad.AdSenseFeedsPos()
	list := make(ASlice, 0)
	for k, v := range data {
		if adPos != 0 && adPos == k {
			list = append(list, &Article{Type: consts.ArticleTypeAdsense})
		}
		list = append(list, v)
	}
	return total, list, err
}
