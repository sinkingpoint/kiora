package raft

import (
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	zerolog.Logger
}

var _ hclog.Logger = (*zerologLogger)(nil)

func (l *zerologLogger) GetLevel() hclog.Level {
	switch l.Logger.GetLevel() {
	case zerolog.TraceLevel:
		return hclog.Trace
	case zerolog.DebugLevel:
		return hclog.Debug
	case zerolog.InfoLevel:
		return hclog.Info
	case zerolog.WarnLevel:
		return hclog.Warn
	case zerolog.ErrorLevel:
		return hclog.Error
	default:
		return hclog.DefaultLevel
	}
}

func (l *zerologLogger) Trace(format string, args ...interface{}) {
	l.Logger.Trace().Fields(args).Msg(format)
}

func (l *zerologLogger) Debug(format string, args ...interface{}) {
	l.Logger.Debug().Fields(args).Msg(format)
}

func (l *zerologLogger) Info(format string, args ...interface{}) {
	l.Logger.Info().Fields(args).Msg(format)
}

func (l *zerologLogger) Warn(format string, args ...interface{}) {
	l.Logger.Warn().Fields(args).Msg(format)
}

func (l *zerologLogger) Error(format string, args ...interface{}) {
	l.Logger.Error().Fields(args).Msg(format)
}

func (l *zerologLogger) Log(level hclog.Level, format string, args ...interface{}) {
	switch level {
	case hclog.Trace:
		l.Logger.Trace().Fields(args).Msg(format)
	case hclog.Debug:
		l.Logger.Debug().Fields(args).Msg(format)
	case hclog.Info:
		l.Logger.Info().Fields(args).Msg(format)
	case hclog.Warn:
		l.Logger.Warn().Fields(args).Msg(format)
	case hclog.Error:
		l.Logger.Error().Fields(args).Msg(format)
	default:
		log.Fatalf("unknown level %d", level)
	}
}

func (l *zerologLogger) SetLevel(level hclog.Level) {
	switch level {
	case hclog.Trace:
		l.Logger = l.Logger.Level(zerolog.TraceLevel)
	case hclog.Debug:
		l.Logger = l.Logger.Level(zerolog.DebugLevel)
	case hclog.Info:
		l.Logger = l.Logger.Level(zerolog.InfoLevel)
	case hclog.Warn:
		l.Logger = l.Logger.Level(zerolog.WarnLevel)
	case hclog.Error:
		l.Logger = l.Logger.Level(zerolog.ErrorLevel)
	default:
		log.Fatalf("unknown level %d", level)
	}
}

func (l *zerologLogger) Name() string {
	return ""
}

func (l *zerologLogger) Named(name string) hclog.Logger {
	return &zerologLogger{l.Logger.With().Str("name", name).Logger()}
}

func (l *zerologLogger) ResetNamed(name string) hclog.Logger {
	return &zerologLogger{l.Logger.With().Str("name", name).Logger()}
}

func (l *zerologLogger) With(args ...interface{}) hclog.Logger {
	return &zerologLogger{l.Logger.With().Fields(args).Logger()}
}

func (l *zerologLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(l.Logger, "", 0)
}

func (l *zerologLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return l.Logger
}

func (l *zerologLogger) ImpliedArgs() []interface{} {
	return nil
}

func (l *zerologLogger) IsTrace() bool {
	return l.Logger.GetLevel() == zerolog.TraceLevel
}

func (l *zerologLogger) IsDebug() bool {
	return l.Logger.GetLevel() == zerolog.DebugLevel
}

func (l *zerologLogger) IsInfo() bool {
	return l.Logger.GetLevel() == zerolog.InfoLevel
}

func (l *zerologLogger) IsWarn() bool {
	return l.Logger.GetLevel() == zerolog.WarnLevel
}

func (l *zerologLogger) IsError() bool {
	return l.Logger.GetLevel() == zerolog.ErrorLevel
}
