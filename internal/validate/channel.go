package validate

type MessagesValidator struct {
	ArticleId string `validate:"required" json:"article_id" form:"article_id"`
	Before    int64  `json:"before" form:"before"`
	After     int64  `json:"after" form:"after"`
	Sort      string `json:"sort" form:"sort"`
	Size      int64  `json:"size" form:"size"`
	Page      int64  `json:"page" form:"page"`
}

type ClientMessage struct {
	Content   string `json:"content" form:"content"`
	ArticleId string `json:"article_id" form:"article_id"`
	Type      string `json:"type" form:"type"`
}
