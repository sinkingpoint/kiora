package serf

// Adapted from https://github.com/sylr/rafty/blob/master/logger/zerolog/hclogger.go
// BSD License.

import (
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
)

// Logger is a wrapper around zerolog.Logger that implements hclog.Logger.
type HCLogger struct {
	zerolog.Logger
}

var _ hclog.Logger = (*HCLogger)(nil)

func (l *HCLogger) GetLevel() hclog.Level {
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
		return hclog.NoLevel
	}
}

func (l *HCLogger) IsTrace() bool {
	return l.Logger.GetLevel() == zerolog.TraceLevel
}

func (l *HCLogger) IsDebug() bool {
	return l.Logger.GetLevel() == zerolog.DebugLevel
}

func (l *HCLogger) IsInfo() bool {
	return l.Logger.GetLevel() == zerolog.InfoLevel
}

func (l *HCLogger) IsWarn() bool {
	return l.Logger.GetLevel() == zerolog.WarnLevel
}

func (l *HCLogger) IsError() bool {
	return l.Logger.GetLevel() == zerolog.ErrorLevel
}

func (l *HCLogger) Trace(format string, args ...interface{}) {
	l.Logger.Trace().Fields(args).Msg(format)
}

func (l *HCLogger) Debug(format string, args ...interface{}) {
	l.Logger.Debug().Fields(args).Msg(format)
}

func (l *HCLogger) Info(format string, args ...interface{}) {
	l.Logger.Info().Fields(args).Msg(format)
}

func (l *HCLogger) Warn(format string, args ...interface{}) {
	l.Logger.Warn().Fields(args).Msg(format)
}

func (l *HCLogger) Error(format string, args ...interface{}) {
	l.Logger.Error().Fields(args).Msg(format)
}

func (l *HCLogger) Log(level hclog.Level, format string, args ...interface{}) {
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

func (l *HCLogger) SetLevel(level hclog.Level) {
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

func (l *HCLogger) Name() string {
	return ""
}

func (l *HCLogger) Named(name string) hclog.Logger {
	return &HCLogger{l.Logger.With().Str("name", name).Logger()}
}

func (l *HCLogger) ResetNamed(name string) hclog.Logger {
	return &HCLogger{l.Logger.With().Str("name", name).Logger()}
}

func (l *HCLogger) With(args ...interface{}) hclog.Logger {
	return &HCLogger{l.Logger.With().Fields(args).Logger()}
}

func (l *HCLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(l.Logger, "", 0)
}

func (l *HCLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return l.Logger
}

func (l *HCLogger) ImpliedArgs() []interface{} {
	return nil
}
