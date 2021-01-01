package main

import (
	"fmt"
	"time"

	logr "github.com/sirupsen/logrus"
)

// LogFormatter for formatting CLI logs
type LogFormatter struct {
	ShowDate bool
}

// Format the log entry
func (l *LogFormatter) Format(entry *logr.Entry) ([]byte, error) {
	var level string
	switch entry.Level {
	case logr.ErrorLevel:
		level = "[ERROR] "
	case logr.FatalLevel:
		level = "[FATAL] "
	case logr.DebugLevel:
		level = "[DEBUG] "
	case logr.WarnLevel:
		level = "[WARN] "
	}
	var date string
	if l.ShowDate {
		date = entry.Time.Format(time.RFC3339) + " "
	}
	return []byte(fmt.Sprintf("%s%s%s", date, level, entry.Message)), nil
}
