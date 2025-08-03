// Package app provides function Run to start full application.
package app

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"CryptocoinPrice/config"
	"CryptocoinPrice/internal/app/server"
	"CryptocoinPrice/internal/pkg/database"
	"CryptocoinPrice/internal/pkg/jsonify"
	"CryptocoinPrice/internal/pkg/logger"
	"CryptocoinPrice/internal/pkg/validator"
)

var _ HTTPServer = (*server.Server)(nil)

// HTTP-server interface.
type HTTPServer interface {
	Run()
	WaitForShutdown() error
}

// Run loads app config and starts HTTP-server and price collector.
// This function is blocking.
func Run() error {
	// create config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}

	// setup logger
	logger.InitLogrus(cfg.App.LogLevel, cfg.App.LogFormat)
	// connect to DB
	gormDB, err := database.New(cfg.DB.ConnString,
		database.WithTranslateError(),
		database.WithIgnoreNotFound(),
		database.WithDisableColorful(),
		database.WithLogLevel(cfg.App.LogLevel),
		database.WithLogger(logrus.StandardLogger()),
	)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	// create json (de)serializer
	jsonifier := jsonify.New()

	// create HTTP-server
	srv, err := server.New(cfg, gormDB, validator.New(), jsonifier)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}
	// run HTTP-server and wait for HTTP-server shutdown
	srv.Run()
	if err := srv.WaitForShutdown(); err != nil {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}
