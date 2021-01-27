package outside

import (
	"errors"
	"time"
)

type Outside struct {
	Id       string    `json:"id"`
	Url      string    `json:"url"`
	Title    string    `json:"title"`
	Priority int       `json:"priority"`
	CreateAt time.Time `json:"create_at"`
	Host     string    `json:"host"`
	Status   int       `json:"status"`
	Email    string    `json:"email"`
}

type OutsideSlice []*Outside

type hit struct {
	Source *Outside `json:"_source"`
}

type hits []hit

var ErrorNoMatch = errors.New("no match")
