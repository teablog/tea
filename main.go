package main

import (
	"github.com/teablog/tea/internal/commands"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "tea"
	app.Version = "v0.3.10"
	app.Commands = []cli.Command{
		commands.Start,
		commands.Deploy,
		commands.AdCode,
		commands.GlobalRegion,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
