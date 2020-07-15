package subscribe

import (
	"dyc/internal/db"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

const SubscriberIndex = "subscriber"

var Email _subscriber

type _subscriber struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

func (s *_subscriber) Store(email string) error {
	res, err := db.ES.Index(
		SubscriberIndex,
		strings.NewReader(fmt.Sprintf(`{
		  "email": "%s",
		  "date": "%s",
		}`, email, time.Now())),
	)
	defer res.Body.Close()
	if err != nil {
		panic(errors.Wrap(err, "es 写入错误"))
	}
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		panic(errors.New(string(resp)))
	}
	return nil
}
