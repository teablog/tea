package article

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
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
			Article   *Article `json:"_source"`
			Highlight struct {
				Content []string `json:"content"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

func (*_search) List(ctx *gin.Context) {
	var (
		buf bytes.Buffer
		r   response
	)
	query := map[string]interface{}{
		"_source": []string{"author", "title", "description", "topic", "id", "date", "last_edit_time"},
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  ctx.Query("q"),
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
		logger.Wrapf(err, "article search json encode err")
		helper.ServerErr(ctx)
		return
	}
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if err != nil {
		logger.Wrapf(err, "article search es err")
		helper.ServerErr(ctx)
		return
	}
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		logger.Errorf("es response err: %s", resp)
		helper.ServerErr(ctx)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Wrapf(err, "article search json decoder err")
		helper.ServerErr(ctx)
		return
	}

	list := make(ASlice, 0, len(r.Hits.Hits))
	total := r.Hits.Total.Value
	for _, v := range r.Hits.Hits {
		v.Article.Highlight = v.Highlight.Content
		list = append(list, v.Article)
	}

	helper.Success(ctx, gin.H{"total": total, "data": list})
	return
}
