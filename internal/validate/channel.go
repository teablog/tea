package validate

type ChannelCreateValidator struct {
	Title   string   `validate:"" json:"title"`
	Members []string `validate:"required" json:"members"`
	Type    string   `validate:"required,oneof=private public" json:"type"`
}

type ChannelMessagesValidator struct {
	ArticleId string `validate:"required" json:"article_id" form:"article_id"`
	Before    int64  `json:"before" form:"before"`
	After     int64  `json:"after" form:"after"`
	Sort      string `json:"sort" form:"sort"`
	Size      int64  `json:"size" form:"size"`
}
