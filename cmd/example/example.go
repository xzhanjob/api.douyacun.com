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
	app.Commands = []cli.Command{
		commands.StartCommands,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	// 日志在这里处理
}
