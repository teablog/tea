package validate

type MessagesValidator struct {
	ArticleId string `validate:"required" json:"article_id" form:"article_id"`
	Before    int64  `json:"before" form:"before"`
	After     int64  `json:"after" form:"after"`
	Sort      string `json:"sort" form:"sort"`
	Size      int `json:"size" form:"size"`
	Page      int  `json:"page" form:"page"`
}

