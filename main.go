package main

import (
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "Gimmick interpreter/compiler"
	app.Usage = "./interpreter --file <filename>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "",
			Usage: "Filename of the script",
		},
	}
	app.Action = func(c *cli.Context) {
		file := c.String("file")
		if file == "" {
			log.Println("Please specify a filename")
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("app.Run() error:", err)
	}
}
