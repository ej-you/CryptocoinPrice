// Package logger provides Init function to setup global logrus logger.
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	_textFormat = "text" // text log format tag
	_jsonFormat = "json" // json log format tag

	_logLevelError = "error" // error log level tag
	_logLevelWarn  = "warn"  // warn log level tag
	_logLevelInfo  = "info"  // info log level tag
)

// JSONFormatterUTC is the logrus.JSONFormatter wrapper with time in UTC.
type JSONFormatterUTC struct {
	logrus.JSONFormatter
}

// Format implements logrus.Formatter and sets time to UTC.
func (f *JSONFormatterUTC) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.JSONFormatter.Format(e)
}

// TextFormatterUTC is the logrus.TextFormatter wrapper with time in UTC.
type TextFormatterUTC struct {
	logrus.TextFormatter
}

// Format implements logrus.Formatter and sets time to UTC.
func (f *TextFormatterUTC) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.TextFormatter.Format(e)
}

// InitLogrus sets up main logger for application with level and formatter.
// Accepted values for logLevel: "info", "warn", "error".
// Accepted values for logFormat: "text", "json".
func InitLogrus(logLevel, logFormat string) {
	logrus.SetOutput(os.Stderr)

	// set formatter
	switch logFormat {
	case _textFormat:
		setTextFormatter()
	case _jsonFormat:
		setJSONFormatter()
	}

	// set log level
	switch logLevel {
	case _logLevelInfo:
		logrus.SetLevel(logrus.InfoLevel)
	case _logLevelWarn:
		logrus.SetLevel(logrus.WarnLevel)
	case _logLevelError:
		logrus.SetLevel(logrus.ErrorLevel)
	}
}

// setTextFormatter sets text formatter for logger.
func setTextFormatter() {
	logrus.SetFormatter(&TextFormatterUTC{
		logrus.TextFormatter{FullTimestamp: true},
	})
}

// setJSONFormatter sets json formatter for logger.
func setJSONFormatter() {
	logrus.SetFormatter(&JSONFormatterUTC{})
}
