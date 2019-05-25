package log

import (
	"github.com/rs/zerolog"
	"os"
)

type ZeroLogger struct {
	zerolog.Logger
}

func NewZeroLogger() ZeroLogger {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return ZeroLogger{logger}
}
