package config

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggingConfig struct {
	ConsoleLoggingEnabled bool
	FileLoggingEnabled    bool
	Directory             string
	Filename              string
	MaxSize               int
	MaxBackups            int
	MaxAge                int
}

func (c *LoggingConfig) Configure() *zerolog.Logger {
	var writers []io.Writer

	if c.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}
	if c.FileLoggingEnabled {
		writers = append(writers, c.newRollingFile())
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	return &logger
}

func (c *LoggingConfig) ConfigureWithoutDisplay() *zerolog.Logger {
	logger := zerolog.New(c.newRollingFile()).With().Timestamp().Logger()

	return &logger
}

func (c *LoggingConfig) newRollingFile() io.Writer {
	if err := os.MkdirAll(c.Directory, 0600); err != nil {
		log.Error().Err(err).Str("path", c.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(c.Directory, c.Filename),
		MaxBackups: c.MaxBackups, // files
		MaxSize:    c.MaxSize,    // megabytes
		MaxAge:     c.MaxAge,     // days
	}
}
