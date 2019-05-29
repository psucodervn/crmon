package subscribers

import (
	"bytes"
	"crmon/pkg/crmon"
	"crmon/pkg/log"
	"encoding/json"
	"fmt"
	"net/http"
)

type MattermostSubscriber struct {
	webHookURL string
	logger     log.ZeroLogger
}

func (s *MattermostSubscriber) Init() error {
	s.logger.Info().Msg(s.Name() + " ready to subscribe")
	return nil
}

func (s *MattermostSubscriber) Cleanup() error {
	s.logger.Info().Msg(s.Name() + " cleaned up")
	return nil
}

func (s *MattermostSubscriber) Name() string {
	return "Mattermost"
}

func (s *MattermostSubscriber) OnReceive(event crmon.Event) error {
	var text string
	switch {
	case event.Action == crmon.ActionInsert:
		text = ":tada: An image was published :tada:"
	case event.Action == crmon.ActionDelete:
		if event.Tag != "" && event.Digest == "" {
			text = ":fire: An image was untagged :fire:"
		} else {
			text = ":fire: An image was deleted :fire:"
		}
	default:
		text = ":bell: New image updated :bell:"
	}

	if event.Tag != "" {
		text += fmt.Sprintf("\n- Tag:    %s", event.Tag)
	}
	if event.Digest != "" {
		text += fmt.Sprintf("\n- Digest: %s", event.Digest)
	}

	payload := map[string]string{
		"username": "Container Registry Monitor",
		"text":     text,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = http.Post(s.webHookURL, "application/json", bytes.NewReader(data))
	return err
}

func NewMattermostSubscriber(webHookURL string) *MattermostSubscriber {
	return &MattermostSubscriber{
		webHookURL: webHookURL,
		logger:     log.NewZeroLogger(),
	}
}
