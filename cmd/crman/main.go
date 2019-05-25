package main

import (
	"crman/pkg/log"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var (
	logger  log.ZeroLogger
	version string
	build   string
)

func init() {
	logger = log.NewZeroLogger()
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s, build %s\n", c.App.Name, c.App.Version, build)
	}

	app := cli.NewApp()
	app.Name = "CRMan"
	app.Version = version
	app.Usage = "monitor image updates in Google Cloud Container Registry"
	app.Email = "hungle@vzota.com.vn"

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
