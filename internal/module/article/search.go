package article

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"io/ioutil"
	"time"
)

var Search _search

type _search struct {
	Author       string    `json:"author"`
	Date         time.Time `json:"date"`
	LastEditTime time.Time `json:"last_edit_time"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Topic        string    `json:"topic"`
	Id           string    `json:"id"`
	Highlight    []string  `json:"highlight"`
}

type response struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Article   Article `json:"_source"`
			Highlight struct {
				Content []string `json:"content"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

func (*_search) List(q string) (int64, ASlice, error) {
	var (
		buf bytes.Buffer
		r   response
	)
	query := map[string]interface{}{
		"_source": []string{"author", "title", "description", "topic", "id", "date", "last_edit_time"},
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q,
				"fields": []string{"title.keyword", "author", "keywords", "content", "label"},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrap(err, "json encode 错误"))
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

	list := make(ASlice, 0, len(r.Hits.Hits))
	total := r.Hits.Total.Value
	for _, v := range r.Hits.Hits {
		v.Article.Highlight = v.Highlight.Content
		list = append(list, &v.Article)
	}
	return total, list, nil
}
