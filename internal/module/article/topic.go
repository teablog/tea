package article

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

var Topics _topic

type _topic struct{}

func (*_topic) List(topic string, page int) (total int64, data []interface{}, err error) {
	var (
		buf bytes.Buffer
		r   map[string]interface{}
	)
	data = make([]interface{}, 0)
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

	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if err != nil {
		panic(errors.Wrap(err, "es search错误"))
	}
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		panic(errors.Wrap(err, "json decode 错误"))
	}

	total = int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, v := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data = append(data, v.(map[string]interface{})["_source"])
	}
	return
}
