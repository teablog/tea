package commands

import (
	"dyc/internal/config"
	"dyc/internal/initialize"
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
			EnvVar:   "_DOUYACUN_CONF",
			Required: false,
			Value:    "/data/web/tea.douyacun.com/configs/prod.ini",
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
