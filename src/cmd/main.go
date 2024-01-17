package cmd

import (
	"math/rand"
	"time"

	"github.com/urfave/cli/v2"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	app := cli.NewApp()
	app.Name = "hroost"
	app.UseShortOptionHandling = true
	app.Setup()
}
