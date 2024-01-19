package main

import (
	"hroost/cmd/command"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	app := command.Root()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
