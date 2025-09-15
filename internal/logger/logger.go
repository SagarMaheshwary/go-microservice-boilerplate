package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
}

type ZerologLogger struct {
	log zerolog.Logger
}

func NewZerologLogger(level string, out io.Writer) *ZerologLogger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	if out == nil {
		out = os.Stderr
	}

	l := zerolog.New(out).With().Timestamp().Logger()
	return &ZerologLogger{log: l}
}

func (l *ZerologLogger) Info(msg string, args ...interface{})  { l.log.Info().Msgf(msg, args...) }
func (l *ZerologLogger) Warn(msg string, args ...interface{})  { l.log.Warn().Msgf(msg, args...) }
func (l *ZerologLogger) Debug(msg string, args ...interface{}) { l.log.Debug().Msgf(msg, args...) }
func (l *ZerologLogger) Error(msg string, args ...interface{}) { l.log.Error().Msgf(msg, args...) }
func (l *ZerologLogger) Fatal(msg string, args ...interface{}) { l.log.Fatal().Msgf(msg, args...) }
func (l *ZerologLogger) Panic(msg string, args ...interface{}) { l.log.Panic().Msgf(msg, args...) }
