package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "krss",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "3000",
				Usage:   "port",
				EnvVars: []string{"PORT"},
			},
		},
		Action: startServerAction,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
