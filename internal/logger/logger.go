package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
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

func (l *ZerologLogger) Info(msg string, fields ...Field) {
	e := l.log.Info()
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (l *ZerologLogger) Warn(msg string, fields ...Field) {
	e := l.log.Warn()
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (l *ZerologLogger) Debug(msg string, fields ...Field) {
	e := l.log.Debug()
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (l *ZerologLogger) Error(msg string, fields ...Field) {
	e := l.log.Error()
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (l *ZerologLogger) Fatal(msg string, fields ...Field) {
	e := l.log.Fatal()
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}

func (l *ZerologLogger) Panic(msg string, fields ...Field) {
	e := l.log.Panic()
	for _, f := range fields {
		e.Interface(f.Key, f.Value)
	}
	e.Msg(msg)
}
