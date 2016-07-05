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
	app.Name = "Unnamed language interpreter"
	app.Usage = "./interpreter --file <filename>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: "",
			Usage: "Filename of the script",
		},
	}
	app.Action = func(c *cli.Context) {
		file := cli.StringFlag("file")
		if file == "" {
			log.Println("Please specify a filename")
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("app.Run() error:", err)
	}
}
