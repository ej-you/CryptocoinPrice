// Package middleware provides all middlewares for HTTP-server.
package middleware

import (
	"os"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// JSON-format for logs
const (
	_textFormat   = "text" // text log format tag
	_jsonFormat   = "json" // json log format tag
	_logLevelInfo = "info" // info log level tag

	_textLogFormat = `INFO[${time}] ${status} | ${method} | ${path} | ${latency}`
	_jsonLogFormat = `{"time": "${time}" "level": "info", "status": "${status}", "method": "${method}", "path": "${path}", "latency": "${latency}"` // nolint:lll // output format
)

// Logger is a middleware for logging all request-response chains.
// Accepted values for logLevel: "info", "warn", "error".
// Accepted values for logFormat: "text", "json".
func Logger(logLevel, logFormat string) fiber.Handler {
	if logLevel != _logLevelInfo {
		return nil
	}

	var format string
	var disableColor bool
	// set formatter
	switch logFormat {
	case _textFormat:
		format = _textLogFormat
		disableColor = false
	case _jsonFormat:
		format = _jsonLogFormat
		disableColor = true
	}

	return logger.New(logger.Config{
		TimeFormat:    time.RFC3339,
		TimeZone:      "UTC",
		Format:        format + "\n",
		Output:        os.Stderr,
		DisableColors: disableColor,
	})
}
