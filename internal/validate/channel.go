package validate

type ChannelCreateValidator struct {
	Title   string   `validate:"" json:"title"`
	Members []string `validate:"required" json:"members"`
	Type    string   `validate:"required,oneof=private public" json:"type"`
}

type ChannelMessagesValidator struct {
	ArticleId string `validate:"required" json:"article_id" form:"article_id"`
	Before    int64  `validate:"required" json:"before" form:"before"`
}
