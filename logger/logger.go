package logger

import (
	"github.com/phuslu/log"
)

var (
	Version = "0.1.0"
)

func GetVersion() string {
	return Version
}

func NewLogger(opts ...OptionsFunc) *log.Logger {
	logger := &log.DefaultLogger

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func NewMultiLogger(opts ...OptionsFunc) *log.Logger {
	logger := &log.DefaultLogger

	logger.Writer = &log.MultiWriter{
		InfoWriter:  &log.FileWriter{LocalTime: true},
		WarnWriter:  &log.FileWriter{LocalTime: true},
		ErrorWriter: &log.FileWriter{LocalTime: true},
		ConsoleWriter: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
		ConsoleLevel: log.DebugLevel,
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}
