package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Uploader",
		Usage: "Resilient and concurrent upload utility",
		Commands: []*cli.Command{
			{
				Name:    "job",
				Aliases: []string{"j"},
				Usage:   "List and run upload jobs",
				Action: func(cCtx *cli.Context) error {
					cfg, _ := LoadConfig("./")
					fmt.Println("DBPath: ", cfg.dbpath)
					return nil
				},
			},
			{
				Name:    "transaction",
				Aliases: []string{"t"},
				Usage:   "List transactions",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
		}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
