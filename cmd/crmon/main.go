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
	date    string

	projectID       string
	topic           string
	subscription    string
	slackWebHookURL string
	command         string
)

func init() {
	logger = log.NewZeroLogger()
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s, build %s at %s\n", c.App.Name, c.App.Version, build, date)
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
		cli.StringFlag{Name: "subscription, s", Destination: &subscription},
		cli.StringFlag{Name: "slack-web-hook, w", Destination: &slackWebHookURL},
		cli.StringFlag{Name: "command, c", Destination: &command},
	}

	app.Action = func(c *cli.Context) error {
		if projectID == "" {
			return errors.New("projectID was required, but not provided")
		}
		if subscription == "" {
			return errors.New("subscription was required, but not provided")
		}

		subs := []crmon.Subscriber{
			subscribers.NewConsoleSubscriber(),
		}
		if slackWebHookURL != "" {
			// subs = append(subs, subscribers.NewSlackSubscriber(slackWebHookURL))
		}
		if command != "" {
			subs = append(subs, subscribers.NewCommandSubscriber(command))
		}
		crApp := crmon.NewApp(crmon.Options{
			ProjectID:    projectID,
			Topic:        topic,
			Subscription: subscription,
			Subscribers:  subs,
		})
		return crApp.Run()
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Err(err).Msg(err.Error())
	}
}
