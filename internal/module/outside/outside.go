package outside

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

func search(body string) (int, OutsideSlice, error) {
	resp, err := db.ES.Search(
		db.ES.Search.WithIndex(OutsideIndex),
		db.ES.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "request err ")
	}
	defer resp.Body.Close()
	if resp.IsError() {
		data, _ := ioutil.ReadAll(resp.Body)
		return 0, nil, errors.Errorf("es response err: %s", string(data))
	}
	r := new(db.ESListResponse)
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return 0, nil, errors.Wrapf(err, "json decode ")
	}
	h := make(hits, 0)
	if err := json.Unmarshal(r.Hits.Hits, &h); err != nil {
		return 0, nil, errors.Wrapf(err, "json unmarshal ")
	}
	data := make(OutsideSlice, 0, len(h))
	for _, v := range h {
		data = append(data, v.Source)
	}
	return r.Hits.Total.Value, data, nil
}

func (row *Outside) create() error {
	row.Id = helper.Md532([]byte(row.Url))
	row.CreateAt = time.Now()
	up, err := url.Parse(row.Url)
	if err != nil {
		return err
	}
	hosts := strings.Split(up.Host, ".")
	if len(hosts) >= 2 {
		row.Host = strings.Join(hosts[len(hosts)-2:], ",")
	} else {
		row.Host = up.Host
	}
	data, err := json.Marshal(row)
	if err != nil {
		return errors.Wrapf(err, "json marshal ")
	}
	resp, err := db.ES.Create(
		OutsideIndex,
		row.Id,
		bytes.NewReader(data),
	)
	if err != nil {
		return errors.Wrapf(err, "request err ")
	}
	defer resp.Body.Close()
	if resp.IsError() {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("es response err: %s", string(b))
	}
	return nil
}

func All() (OutsideSlice, error) {
	query := `{"size":10000,"_source":"host"}`
	_, list, err := search(query)
	if err != nil {
		return nil, err
	}
	return list, nil
}
