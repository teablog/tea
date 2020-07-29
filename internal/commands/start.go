package commands

import (
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/initialize"
	"github.com/urfave/cli"
)

var Start = cli.Command{
	Name:   "start",
	Usage:  "",
	Action: startAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "conf",
			Usage:    "-conf <path>",
			EnvVar:   "_TEA_CONF",
			Required: false,
		},
	},
}

func startAction(c *cli.Context) (err error) {
	// 加载配置文件
	config.Init(c.String("conf"))
	// 启动web服务
	initialize.Server()

	return nil
}
