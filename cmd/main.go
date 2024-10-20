package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jiazhen-lin/flight-booking/internal/adapter"
	"github.com/jiazhen-lin/flight-booking/internal/app"
	"github.com/jiazhen-lin/flight-booking/internal/service"
)

type config struct {
	httpAddr string
	dbArg    string
}

func initConfig() (*config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("empty env HTTP_ADDR")
	}
	dbArg := os.Getenv("DB_ARG")
	if dbArg == "" {
		return nil, fmt.Errorf("empty env DB_ARG")
	}

	return &config{
		httpAddr: addr,
		dbArg:    dbArg,
	}, nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	// init
	config, err := initConfig()
	if err != nil {
		logrus.Fatalf("init config: %v", err)
	}

	engine, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: config.dbArg,
			},
		),
	)
	if err != nil {
		logrus.Fatalf("gorm.Open error: %v", err)
	}

	db, err := engine.DB()
	if err != nil {
		logrus.Fatalf("get db error: %v", err)
	}
	db.SetConnMaxLifetime(time.Duration(10) * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	flightRepo := adapter.NewFlightPostgresRepository(engine)
	flightService := service.NewFlightService(flightRepo)
	application := app.NewApplication(flightService)
	server := app.NewHTTPServer(config.httpAddr, application)

	// start http server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Errorf("http server ListenAndServe error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logrus.Info("received stop signal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("http server shutdown error: %v", err)
	}

	logrus.Info("http server shutdown gracefully")
}
