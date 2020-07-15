package validate

type ChannelCreateValidator struct {
	Title   string   `validate:"" json:"title"`
	Members []string `validate:"required" json:"members"`
	Type    string   `validate:"required,oneof=private public" json:"type"`
}

type ChannelMessagesValidator struct {
	ChannelId string `validate:"required" json:"channel_id" form:"channel_id"`
	Before    int64 `validate:"required" json:"before" form:"before"`
}
