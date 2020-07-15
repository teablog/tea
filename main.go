package main

import (
	"dyc/internal/commands"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "douyacun"
	app.Version = "v0.3.10"
	app.Commands = []cli.Command{
		commands.Start,
		commands.Deploy,
		commands.AdCode,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
