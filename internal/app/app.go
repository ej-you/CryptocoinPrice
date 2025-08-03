// Package app provides struct with Run method to start full application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"

	"CryptocoinPrice/config"
	"CryptocoinPrice/internal/app/pricecollector"
	"CryptocoinPrice/internal/app/server"
	"CryptocoinPrice/internal/pkg/database"
	"CryptocoinPrice/internal/pkg/jsonify"
	"CryptocoinPrice/internal/pkg/logger"
	"CryptocoinPrice/internal/pkg/validator"
)

var _ Service = (*pricecollector.PriceCollector)(nil)

// App service interface.
type Service interface {
	StartWithShutdown(ctx context.Context) error
}

type App struct {
	cfg      *config.Config
	services []Service
}

// New returns new app instance.
func New() (*App, error) {
	// load config
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("create config: %w", err)
	}
	// setup logger
	logger.InitLogrus(cfg.App.LogLevel, cfg.App.LogFormat)

	// connect to DB
	gormDB, err := database.New(cfg.DB.ConnString,
		database.WithTranslateError(),
		database.WithIgnoreNotFound(),
		database.WithDisableColorful(),
		database.WithLogLevel(cfg.App.LogLevel),
		database.WithLogger(logrus.StandardLogger()))
	if err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	// init serv
	srv, err := server.New(cfg, gormDB, validator.New(), jsonify.New())
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}
	// init price collector
	priceCollector := pricecollector.New(cfg, gormDB)

	return &App{
		cfg:      cfg,
		services: []Service{srv, priceCollector},
	}, nil
}

// Run starts HTTP-server service and price collector service.
// This function is blocking. It waits for os signal to gracefully shutdown all services.
func (a *App) Run() error {
	var appErr error

	// ctx for app
	appContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// handle shutdown process signals
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// start all services
	var wg sync.WaitGroup // nolint:varnamelen // generally accepted name
	serviceErr := make(chan error, 1)
	for _, service := range a.services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := service.StartWithShutdown(appContext); err != nil {
				serviceErr <- err
			}
		}()
	}

	select {
	case handledSignal := <-quitSig:
		cancel()
		logrus.Infof("Got %s signal. Shutdown services...", handledSignal.String())
	case err := <-serviceErr:
		cancel()
		appErr = fmt.Errorf("service: %w", err)
		logrus.Info("One of the services fell down. Shutdown other services...")
	case <-appContext.Done():
		appErr = appContext.Err()
		logrus.Info("Context canceled. Shutdown app...")
	}

	// wait for all services
	wg.Wait()
	logrus.Info("All services was stopped. Shutdown app")
	return appErr
}
