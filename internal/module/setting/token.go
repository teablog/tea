package setting

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
)

var Token *token

type token struct{}

func (*token) Get(ctx *gin.Context) string {
	return fmt.Sprintf("**token**:\n\n`%s`\n\n**限流**:\n\n- `50次/秒`\n\n- `500000次/天`\n\n> 如果需要更高频的调用，请联系 douyacun@gmail.com ，感谢您的使用！\n\n", config.GetKey("help::token").String())
}
