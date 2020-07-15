package commands

import (
	"dyc/internal/config"
	"dyc/internal/db"
	"dyc/internal/logger"
	"dyc/internal/module/deploy"
	"github.com/urfave/cli"
)

var Deploy = cli.Command{
	Name: "deploy",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "conf",
			Usage:    "-conf <path>",
			Required: false,
			Value:    "configs/prod.ini",
		},
		cli.StringFlag{
			Name:     "dir",
			Usage:    "指定文章所在目录",
			Required: false,
			Value:    "/data/book/",
		},
	},
	Action: deployAction,
}

func deployAction(c *cli.Context) (err error) {
	// 加载配置文件
	config.Init(c.String("conf"))
	// 设置运行环境
	logger.NewLogger(config.GetLogFD())
	logger.SetLevel("debug")
	// 数据库
	db.NewElasticsearch(config.GetKey("elasticsearch::address").Strings(","), config.GetKey("elasticsearch::user").String(), config.GetKey("elasticsearch::password").String())
	deploy.Run(c.String("dir"))
	return
}
