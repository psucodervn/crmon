package subscribers

import (
	"crmon/pkg/crmon"
	"crmon/pkg/log"
	"errors"
	"os/exec"
	"strings"
)

type commandSubscriber struct {
	logger  log.ZeroLogger
	command string
	args    []string
}

func (s *commandSubscriber) Name() string {
	return "Command"
}

func (s *commandSubscriber) Init() error {
	exists := s.checkCommandExists(s.command)
	if !exists {
		return errors.New("command " + s.command + " is not executable")
	}
	s.logger.Info().Msg(s.Name() + " ready to subscribe")
	return nil
}

func (s *commandSubscriber) Cleanup() error {
	s.logger.Info().Msg(s.Name() + " cleaned up")
	return nil
}

func (s *commandSubscriber) OnReceive(event crmon.Event) error {
	cmd := exec.Command(s.command, s.args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Error().Err(err).Msg(string(out))
		return err
	}
	s.logger.Info().Msg(string(out))
	return nil
}

func (s *commandSubscriber) checkCommandExists(command string) bool {
	path, err := exec.LookPath(command)
	if err != nil {
		s.logger.Warn().Str("command", command).Msg("not found")
		return false
	} else {
		s.logger.Debug().Str("command", command).Str("path", path).Msg("found")
		return true
	}
}

func NewCommandSubscriber(command string) crmon.Subscriber {
	ss := strings.Split(command, " ")
	s := &commandSubscriber{
		logger:  log.NewZeroLogger(),
		command: ss[0],
		args:    ss[1:],
	}
	return s
}
