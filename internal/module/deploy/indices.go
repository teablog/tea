package deploy

import (
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/logger"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var (
	Indices _Indices
)

type _article struct{}

type _Indices struct {
	Article _article
}

func (*_article) Create(index string) error {
	res, err := db.ES.Indices.Exists(
		[]string{index},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		res2, err := db.ES.Indices.Create(
			index,
			db.ES.Indices.Create.WithBody(strings.NewReader(consts.IndicesTopicMapping)),
		)
		if err != nil {
			return err
		}
		if res2.IsError() {
			resp, _ := ioutil.ReadAll(res.Body)
			return errors.New(string(resp))
		}
	}
	return nil
}

func (*_article) Delete(index string) error {
	res, err := db.ES.Indices.Exists(
		[]string{index},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		res2, err := db.ES.Indices.Delete(
			[]string{index},
		)
		if err != nil {
			return err
		}
		if res2.IsError() {
			resp, _ := ioutil.ReadAll(res.Body)
			return errors.New(string(resp))
		}
	}
	return nil
}

func (*_article) ReindexAndDeleteSource(source, dest string) (err error) {
	logger.Debugf("create indices(%s) mapping", dest)
	if err = Indices.Article.Create(dest); err != nil {
		return
	}
	query := fmt.Sprintf(`{
	  "source": {
		"index": "%s"
	  },
	  "dest": {
		"index": "%s"
	  }
	}`, source, dest)
	logger.Debugf("reindex source(%s) desc(%s)", source, dest)
	res, err := db.ES.Reindex(
		strings.NewReader(query),
		db.ES.Reindex.WithRefresh(true),
	)
	if err != nil {
		return errors.Wrap(err, "es reindex错误")
	}
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	if res.IsError() {
		return errors.New(string(resp))
	}
	logger.Debugf("reindex response: %s", resp)
	return Indices.Article.Delete(source)
}
