package logger_infrastructure

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/Aldiwildan77/backend-hexa-template/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once

func NewLogger(conf config.Config) zerolog.Logger {
	var log zerolog.Logger

	once.Do(func() {
		logWriters := make([]io.Writer, 0)

		if conf.Logger.ConsoleLogEnabled {
			logWriters = append(logWriters, NewConsoleLogger(ConsoleLogger{
				Out:        os.Stderr,
				NoColor:    conf.Logger.ConsoleLogNoColor,
				TimeFormat: time.RFC3339,
			}))
		}

		if conf.Logger.FileLogEnabled {
			logWriters = append(logWriters, NewLumberjackFileLogger(LumberjackFileLogger{
				Directory:  conf.Logger.Directory,
				Filename:   conf.Logger.Filename,
				MaxBackups: conf.Logger.MaxBackups,
				MaxSize:    conf.Logger.MaxSize,
				MaxAge:     conf.Logger.MaxAge,
				LocalTime:  conf.Logger.LocalTime,
				Compress:   conf.Logger.Compress,
			}))
		}

		zerolog.TimestampFunc = time.Now
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		mw := zerolog.MultiLevelWriter(logWriters...)

		log = zerolog.
			New(mw).
			With().
			Timestamp().
			Logger()
	})

	return log
}
