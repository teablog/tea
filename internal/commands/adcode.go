package commands

import (
	"bytes"
	"dyc/internal/config"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/logger"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"strings"
)

var AdCode = cli.Command{
	Name:        "adcode",
	Description: "获胜所属省级市级code",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "csv",
			Usage:    "--csv [path]",
			Required: true,
		},
		cli.StringFlag{
			Name:     "env",
			Usage:    "--env [env]",
			Required: false,
			Value:    "debug",
		},
	},
	Action: adcode,
}
// 城市code：https://lbs.amap.com/api/webservice/download
func adcode(c *cli.Context) error {
	config.Init(c.String("env"))
	db.NewElasticsearch(config.GetKey("elasticsearch::address").Strings(","), config.GetKey("elasticsearch::user").String(), config.GetKey("elasticsearch::password").String())
	logger.NewLogger(config.GetLogFD())
	fp, err := os.Open(c.String("csv"))
	if err != nil {
		return err
	}
	r := csv.NewReader(fp)
	if recodes, err := r.ReadAll(); err != nil {
		return err
	} else {
		truncateIndices()
		store(recodes)
	}
	return nil
}

func truncateIndices() error {
	res, err := db.ES.Indices.Exists(
		[]string{consts.IndicesAdCodeConst},
	)
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		res, err := db.ES.Indices.Delete(
			[]string{consts.IndicesAdCodeConst},
		)
		if err != nil {
			return err
		}
		esResponsePrint(res, false)
	}
	if res, err := db.ES.Indices.Create(
		consts.IndicesAdCodeConst,
		db.ES.Indices.Create.WithBody(strings.NewReader(consts.IndicesAdCodeMapping)),
	); err != nil {
		return err
	} else {
		esResponsePrint(res, false)
	}
	logger.Debugf("%s recreate success", consts.IndicesAdCodeConst)
	return nil
}

func esResponsePrint(response *esapi.Response, anyway bool) {
	body, _ := ioutil.ReadAll(response.Body)
	if response.IsError() {
		logger.Errorf("es response error: %s", body)
		return
	}
	if anyway {
		logger.Debugf("es response: %s", body)
	}
}

func store(records [][]string) error {
	header := records[0]
	start := 0
	name, adCode, cityCode := 0, 0, 0
	for k, v := range header {
		switch v {
		case "adcode":
			adCode = k
			start = 1
		case "citycode":
			cityCode = k
		case "中文名":
			name = k
		}
	}
	buf := bytes.NewBuffer(nil)
	indexBulk := fmt.Sprintf(`{ "index":{"_index":"%s"}}`, consts.IndicesAdCodeConst)
	for i := start; i < len(records); i++ {
		buf.WriteString(indexBulk)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf(`{"name": "%s", "adcode": "%s", "citycode": "%s"}`, records[i][name], records[i][adCode], records[i][cityCode]))
		buf.WriteString("\n")
	}
	//logger.Debugf("es bulk body: %s", buf.String())
	res, err := db.ES.Bulk(buf)
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.Wrap(err, "es response err")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	type _esResponse struct {
		Took int `json:"took"`
		Errors bool `json:"errors"`
	}
	var esResponse _esResponse
	if err = json.Unmarshal(body, &esResponse); err != nil {
		return err
	}
	if esResponse.Errors {
		return errors.New("bulk exec failed!!!!!")
	}
	logger.Debugf("down total: %d", len(records))
	return nil
}
