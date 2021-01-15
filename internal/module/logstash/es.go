package logstash

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"io/ioutil"
	"strings"
	"time"
)

var ES *_es

type _es struct{}

func (*_es) KongHttpLog(data string) error {
	index := fmt.Sprintf(consts.SpiderIndices, time.Now().Format("200601"))
	resp, err := db.ES.Index(
		index,
		strings.NewReader(data),
	)
	if err != nil {
		return errors.Wrap(err, "fetch es err")
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read es response err")
	}
	if resp.IsError() {
		return errors.Errorf("es response err: %s", string(rb))
	}
	return nil
}
