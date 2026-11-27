package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Entry for structured logging.
type Logger struct {
	*logrus.Entry
}

// NewLogger creates a new logger with the given log level.
func NewLogger(level string) (*Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	base := logrus.New()
	base.SetOutput(os.Stdout)
	base.SetLevel(lvl)
	base.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return &Logger{Entry: logrus.NewEntry(base)}, nil
}