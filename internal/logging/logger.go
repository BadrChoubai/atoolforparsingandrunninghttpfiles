package logging

import (
	"io"
	"log/slog"
)

type StructuredLogger interface {
	Info(msg string, args ...any)
	Error(whatWasHappening string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Logger struct {
	stdout *slog.Logger
	stderr *slog.Logger
}

func NewLogger(stdout, stderr io.Writer) *Logger {
	return &Logger{
		stdout: slog.New(slog.NewJSONHandler(stdout, nil)),
		stderr: slog.New(slog.NewJSONHandler(stderr, nil)),
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.stdout.Info(msg, args...)
}

func (l *Logger) Error(whatWasHappening string, args ...any) {
	l.stderr.Error(whatWasHappening, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.stdout.Warn(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.stdout.Debug(msg, args...)
}
