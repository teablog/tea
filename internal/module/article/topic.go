package article

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"strconv"
)

var Topics _topic

type _topic struct{}

func (*_topic) List(ctx *gin.Context){
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	var (
		buf bytes.Buffer
	)
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
					"value": ctx.Query("topic"),
				},
			},
		},
		"_source": []string{"author", "title", "description", "topic", "id", "cover", "date", "last_edit_time"},
	}
	if err = json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Wrapf(err, "topic search json encode err")
		helper.ServerErr(ctx)
		return
	}
	total, data, err := Art.Search(buf.String())
	if err != nil {
		logger.Wrapf(err, "topic search err")
		helper.ServerErr(ctx)
		return
	}
	adPos := config.Ad.AdSenseFeedsPos()
	list := make(ASlice, 0)
	for k, v := range data {
		if adPos != 0 && adPos == k {
			list = append(list, &Article{Type: consts.ArticleTypeAdsense})
		}
		list = append(list, v)
	}
	helper.Success(ctx, gin.H{"total": total, "data": data})
	return
}
