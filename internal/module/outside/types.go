package outside

import "time"

type Outside struct {
	Id       string    `json:"id"`
	Url      string    `json:"url"`
	Title    string    `json:"title"`
	Priority int       `json:"priority"`
	CreateAt time.Time `json:"create_at"`
	Status   int       `json:"status"`
}

type OutsideSlice []*Outside

type hit struct {
	Source *Outside `json:"_source"`
}

type hits []hit
