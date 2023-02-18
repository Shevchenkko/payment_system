// Package provides a logging functionality.
package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Interface - represents logger interface.
type Interface interface {
	Named(name string) Interface
	With(args ...interface{}) Interface
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

// Logger - represents instance of logger.
type Logger struct {
	logger *zap.SugaredLogger
}

var _ Interface = (*Logger)(nil)

// New - creates new instance logger.
func New(level string) *Logger {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	// logger config
	config := zap.Config{
		Development:      false,
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(l),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeDuration: zapcore.SecondsDurationEncoder,
			LevelKey:       "severity",
			EncodeLevel:    zapcore.CapitalLevelEncoder, // e.g. "Info"
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder, // e.g. package/file:line
			TimeKey:        "timestamp",
			EncodeTime:     zapcore.ISO8601TimeEncoder, // e.g. 2020-05-05T03:24:36.903+0300
			NameKey:        "name",
			EncodeName:     zapcore.FullNameEncoder, // e.g. GetSiteGeneralHandler
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     "\n",
		},
	}

	// build logger from config
	logger, _ := config.Build()

	// configure and create logger
	return &Logger{
		logger: logger.Sugar(),
	}
}

// Named - returns a new logger with a chained name.
func (l *Logger) Named(name string) Interface {
	return &Logger{
		logger: l.logger.Named(name),
	}
}

// Named - returns a new logger with a chained name.
func (l *Logger) With(args ...interface{}) Interface {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// Debug - logs in debug level.
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debugw(message, args...)
}

// Info - logs in info level.
func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Infow(message, args...)
}

// Warn - logs in warn level.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warnw(message, args...)
}

// Error - logs in error level.
func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Errorw(message, args...)
}

// Fatal - logs and exits program with status 1.
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalw(message, args...)
	os.Exit(1)
}
