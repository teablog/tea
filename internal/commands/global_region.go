package commands

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

var GlobalRegion = cli.Command{
	Name:        "region",
	Description: "import global region from csv",
	Action:      globalRegion,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "csv",
			Usage:    "--csv [path]",
			EnvVar:   "__TEA_REGION_CSV",
			FilePath: "",
			Required: true,
		},
		cli.StringFlag{
			Name:     "conf",
			Usage:    "--conf [path]",
			EnvVar:   "__TEA_CONF",
			FilePath: "",
			Required: true,
		},
	},
}

func globalRegion(ctx *cli.Context) error {
	config.Init(ctx.String("conf"))
	db.NewElasticsearch(config.GetKey("elasticsearch::address").Strings(","), config.GetKey("elasticsearch::user").String(), config.GetKey("elasticsearch::password").String())

	file, err := os.Open(ctx.String("csv"))
	if err != nil {
		return errors.Wrap(err, "csv file open failed:")
	}
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return errors.Wrap(err, "csv file read failed:")
	}
	buf := bytes.NewBuffer(nil)
	for _, record := range records {
		buf.WriteString(fmt.Sprintf(`{ "index" : { "_index" : "%s", "_id" : "%s" } }`, consts.IndicesGlobalRegion, record[0]))
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf(`{"pid":%s,"path":"%s","level":%s,"name":"%s","name_en":"%s","name_pinyin":"%s","code":"%s"}`,
			record[1], record[2], record[3], record[4], record[5], record[6], record[7]))
		buf.WriteString("\n")
	}
	resp, err := db.ES.Bulk(buf)
	if err != nil {
		return errors.Wrap(err, "es bulk execute failed")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.IsError() {
		return errors.Errorf("es error: %s", body)
	}
	fmt.Printf("down, total: %d", len(records))
	return nil
}
