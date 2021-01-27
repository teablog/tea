package outside

import (
	"errors"
	"github.com/teablog/tea/internal/helper"
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

func NewOutside() *Outside {
	return &Outside{}
}

func (row *Outside) GenId() string {
	if row.Url != "" {
		row.Id = helper.Md532([]byte(row.Url))
		return row.Id
	}
	return ""
}

func (row *Outside) SetUrl(url string) *Outside {
	row.Url = url
	return row
}

func (row *Outside) SetTitle(title string) *Outside {
	row.Title = title
	return row
}

func (row *Outside) SetEmail(email string) *Outside {
	row.Email = email
	return row
}

type OutsideSlice []*Outside

type hit struct {
	Source *Outside `json:"_source"`
}

type hits []hit

var ErrorNoMatch = errors.New("no match")
