package outside

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
)

func Route(api *gin.RouterGroup) {
	o := api.Group("/outside")
	{
		o.GET("", list)
		o.POST("", post)
	}
}

func list(ctx *gin.Context) {
	body := fmt.Sprintf(`
{
    "query": {
        "term": {
            "status": %d
        }
    },
    "sort": [
        {
            "priority": {
                "order": "desc"
            }
        },
        {
            "create_at": {
                "order": "asc"
            }
        }
    ],
    "size": 20
}
`, consts.StatusOn)
	c, data, err := search(body)
	if err != nil {
		logger.Wrapf(err, "outside es search err ")
		helper.ServerErr(ctx)
		return
	}
	helper.Success(ctx, gin.H{"total": c, "list": data})
	return
}

func post(ctx *gin.Context) {
	o := new(Outside)
	if err := ctx.Bind(o); err != nil {
		helper.Fail(ctx, errors.New("参数错误"))
		return
	}
	if o.Url == "" || o.Title == "" {
		helper.Fail(ctx, errors.New("url/title 不能为空～"))
		return
	}
	if err := o.create(); err != nil {
		logger.Wrapf(err, "outside es create err ")
		helper.ServerErr(ctx)
		return
	}
	helper.Success(ctx, o)
	return
}
