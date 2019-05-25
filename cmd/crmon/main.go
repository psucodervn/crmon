package main

import (
	"crmon/pkg/crmon"
	"crmon/pkg/crmon/subscribers"
	"crmon/pkg/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
)

var (
	logger  log.ZeroLogger
	version string
	build   string

	projectID    string
	topic        string
	subscription string
)

func init() {
	logger = log.NewZeroLogger()
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s, build %s\n", c.App.Name, c.App.Version, build)
	}

	app := cli.NewApp()
	app.Name = "Container Registry Monitor"
	app.HelpName = "crmon"
	app.Version = version
	app.Usage = "monitor image updates in Google Cloud Container Registry"
	app.Email = "hungle@vzota.com.vn"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "project-id, p", Destination: &projectID},
		cli.StringFlag{Name: "topic, t", Destination: &topic, Value: "gcr"},
		cli.StringFlag{Name: "subscription, s", Destination: &subscription, Value: "gcr"},
	}

	app.Action = func(c *cli.Context) error {
		if projectID == "" {
			return errors.New("projectID was required, but not provided")
		}
		crApp := crmon.NewApp(crmon.Options{
			ProjectID:    projectID,
			Topic:        topic,
			Subscription: subscription,
			Subscribers: []crmon.Subscriber{
				subscribers.NewConsoleSubscriber(),
			},
		})
		return crApp.Run()
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Err(err).Msg(err.Error())
	}
}
