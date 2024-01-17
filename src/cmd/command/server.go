package command

import (
	"hroost/server"

	"github.com/urfave/cli/v2"
)

func serverCommand() *cli.Command {
	return &cli.Command{

		Name:     "server",
		Usage:    "run the application server",
		Category: "Server",
		Action: func(cCtx *cli.Context) error {
			server := server.NewServer()
			err := server.Run(cCtx.Context)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
