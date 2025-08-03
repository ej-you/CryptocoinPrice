// Package server provides HTTP-server interface.
package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	cfg     *config.Config
	db      *gorm.DB
	valid   validator.Validator
	jsonify jsonify.Jsonify

	fiberApp *fiber.App
	err      chan error // server listen error
}

// New returns new server instance.
func New(cfg *config.Config, db *gorm.DB,
	valid validator.Validator, jsonifier jsonify.Jsonify) (*Server, error) {

	return &Server{
		cfg:     cfg,
		db:      db,
		valid:   valid,
		jsonify: jsonifier,
		err:     make(chan error),
	}, nil
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
// Run starts server.
func (s *Server) Run() {
	// app init
	s.fiberApp = fiber.New(fiber.Config{
		JSONEncoder:   s.jsonify.Marshal,
		JSONDecoder:   s.jsonify.Unmarshal,
		ServerHeader:  "Cryptocoin Price API",
		StrictRouting: false,
	})

	// set up base middlewares
	httpLogger := middleware.Logger(s.cfg.App.LogLevel, s.cfg.App.LogFormat)
	if httpLogger != nil {
		s.fiberApp.Use(httpLogger)
	}
	s.fiberApp.Use(middleware.Recover())
	s.fiberApp.Use(middleware.Swagger())
	// register all endpoints
	s.registerEndpointsV1()

	// start app
	go func() {
		if err := s.fiberApp.Listen(":" + s.cfg.Server.Port); err != nil {
			s.err <- fmt.Errorf("listen: %w", err)
		}
	}()
}

// WaitForShutdown waits for OS signal to gracefully shuts down server.
// This method is blocking.
func (s *Server) WaitForShutdown() error {
	// skip if server is not running
	if s.fiberApp == nil {
		return nil
	}

	// handle shutdown process signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	shutdownDone := make(chan struct{})
	// create gracefully shutdown task
	var err error
	go func() {
		defer close(shutdownDone)
		select {
		case err = <-s.err: // server listen error
			return
		case handledSignal := <-quit:
			logrus.Infof("Got %s signal. Shutdown server", handledSignal.String())
			// shutdown app
			s.fiberApp.ShutdownWithTimeout(s.cfg.Server.ShutdownTimeout) // nolint:errcheck // cannot occurs
		}
	}()

	// wait for shutdown
	<-shutdownDone
	logrus.Info("Server shutdown successfully")
	return err
}
