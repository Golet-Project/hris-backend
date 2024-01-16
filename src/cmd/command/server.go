package command

import (
	"github.com/urfave/cli/v2"
)

func serverCommand() *cli.Command {
	return &cli.Command{

		Name:     "server",
		Usage:    "run the application server",
		Category: "Server",
		Action: func(cCtx *cli.Context) error {
			RunApp(cfg)

			return nil
		},
	}
}
