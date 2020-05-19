package main

import (
	"github.com/dairovolzhas/dar-project/game"
	"github.com/urfave/cli"
	"log"
	"math/rand"
	"os"
	"time"
)
var (
	flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "logging",
			Aliases: []string{"l"},
			Usage:   "--logging [-l] to enable logging, default disabled",
			Destination: &game.Logs,
		},
	}
)

func main() {
	app := &cli.App{
		Flags:flags,
		Name: "Lenin's Tank",
		Usage: " 'F1': to go MENU; 'Arrows': to MOVE; 'Space': to fire; ",
		Action: run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


func run(c *cli.Context) error {
	if err := Start(); err != nil {
		return err
	}
	return nil
}

func Start() (err error) {
	if game.Logs {
		file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		log.SetOutput(file)
	}

	rand.Seed(time.Now().UnixNano())

	err = game.RabbitMQ()
	defer game.CloseConnectionAndChannel()

	game.Game().Start()
	return
}