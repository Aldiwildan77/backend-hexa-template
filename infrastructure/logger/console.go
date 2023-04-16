package logger_infrastructure

import (
	"io"

	"github.com/rs/zerolog"
)

type ConsoleLogger struct {
	Out        io.Writer
	NoColor    bool
	TimeFormat string
}

func NewConsoleLogger(conf ConsoleLogger) io.Writer {
	return zerolog.ConsoleWriter{
		Out:        conf.Out,
		NoColor:    conf.NoColor,
		TimeFormat: conf.TimeFormat,
	}
}
