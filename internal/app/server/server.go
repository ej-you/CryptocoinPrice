// Package server provides HTTP-server interface.
package server

import (
	"context"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"CryptocoinPrice/config"
	"CryptocoinPrice/internal/app/server/middleware"

	"CryptocoinPrice/internal/pkg/jsonify"
	"CryptocoinPrice/internal/pkg/validator"
)

// HTTP-server.
type Server struct {
	cfg      *config.Config
	fiberApp *fiber.App
}

//	@title			Cryptocoin Price API
//	@version		1.0.0
//	@description	HTTP API для сбора, хранения и отображения стоимости криптовалют.
//
//	@host			127.0.0.1:8000
//	@basePath		/api/v1
//	@schemes		http
//
//	@accept			json
//	@produce		json
//
// New returns new server instance.
func New(cfg *config.Config, dbStorage *gorm.DB,
	valid validator.Validator, jsonifier jsonify.Jsonify) (*Server, error) {

	// fiber init
	server := &Server{
		cfg: cfg,
		fiberApp: fiber.New(fiber.Config{
			JSONEncoder:   jsonifier.Marshal,
			JSONDecoder:   jsonifier.Unmarshal,
			ServerHeader:  "Cryptocoin Price API",
			StrictRouting: false,
		}),
	}

	// set up base middlewares
	httpLogger := middleware.Logger(cfg.App.LogLevel, cfg.App.LogFormat)
	if httpLogger != nil {
		server.fiberApp.Use(httpLogger)
	}
	server.fiberApp.Use(middleware.Recover())
	server.fiberApp.Use(middleware.Swagger())
	// register all endpoints
	server.registerEndpointsV1(cfg, dbStorage, valid)

	return server, nil
}

// StartWithShutdown starts server and waits for
// context is done for gracefully shutdown server.
// This method is blocking.
func (s *Server) StartWithShutdown(ctx context.Context) error {
	logrus.Info("Start server")
	defer logrus.Info("Server is shutdown")

	errChan := make(chan error, 1)
	defer close(errChan)
	// start server
	go func() {
		if err := s.fiberApp.Listen(":" + s.cfg.Server.Port); err != nil {
			errChan <- fmt.Errorf("server: listen: %w", err)
		}
	}()

	// wait for context or server listen error
	select {
	case <-ctx.Done():
		if err := s.fiberApp.ShutdownWithTimeout(s.cfg.App.ShutdownTimeout); err != nil {
			return fmt.Errorf("server: shutdown: %w", err)
		}
		return nil
	case err := <-errChan:
		return err
	}
}
