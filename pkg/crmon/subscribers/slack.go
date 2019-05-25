package subscribers

import (
	"crmon/pkg/crmon"
	"github.com/ashwanthkumar/slack-go-webhook"
)

type slackSubscriber struct {
	webHookURL string
}

func (s *slackSubscriber) Name() string {
	return "Slack"
}

func (s *slackSubscriber) OnReceive(event crmon.Event) error {
	att := slack.Attachment{}
	att.AddField(slack.Field{Title: "Digest", Value: event.Digest})
	if len(event.Tag) > 0 {
		att.AddField(slack.Field{Title: "Tag", Value: event.Tag})
	}

	payload := slack.Payload{
		Username:    "CR Monitor",
		Attachments: []slack.Attachment{att},
	}
	switch {
	case event.Action == crmon.ActionInsert:
		payload.Text = ":tada: An image was published :tada:"
	case event.Action == crmon.ActionDelete:
		payload.Text = ":fire: An image was deleted :fire:"
	default:
		payload.Text = ":bell: New image updated :bell:"
	}

	errs := slack.Send(s.webHookURL, "", payload)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func NewSlackSubscriber(webHookURL string) crmon.Subscriber {
	return &slackSubscriber{
		webHookURL: webHookURL,
	}
}
