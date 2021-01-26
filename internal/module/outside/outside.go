package outside

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"io/ioutil"
	"strings"
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

func create(row *Outside) error {
	row.Id = helper.Md532([]byte(row.Url))
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
