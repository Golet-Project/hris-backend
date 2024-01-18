package command

import "github.com/urfave/cli/v2"

func Root() *cli.App {
	app := cli.NewApp()
	app.Name = "hroost"
	app.UseShortOptionHandling = true
	app.Commands = []*cli.Command{
		migrationCommand(),
		serverCommand(),
	}

	app.Setup()

	return app
}
