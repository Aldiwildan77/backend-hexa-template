package logger_infrastructure

import (
	"io"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LumberjackFileLogger struct {
	Directory  string
	Filename   string
	MaxBackups int
	MaxSize    int
	MaxAge     int
	LocalTime  bool
	Compress   bool
}

func NewLumberjackFileLogger(conf LumberjackFileLogger) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(conf.Directory, conf.Filename),
		MaxBackups: conf.MaxBackups, // files
		MaxSize:    conf.MaxSize,    // megabytes
		MaxAge:     conf.MaxAge,     // days
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	}
}
