package subscribers

import (
	"crmon/pkg/crmon"
	"crmon/pkg/log"
)

type consoleSubscriber struct {
	logger log.ZeroLogger
}

func NewConsoleSubscriber() crmon.Subscriber {
	return &consoleSubscriber{
		logger: log.NewZeroLogger(),
	}
}

func (s *consoleSubscriber) Name() string {
	return "Console"
}

func (s *consoleSubscriber) OnReceive(event crmon.Event) error {
	s.logger.Info().
		Str("action", event.Action).
		Str("tag", event.Tag).
		Str("digest", event.Digest).
		Msg("new image updates")
	return nil
}
