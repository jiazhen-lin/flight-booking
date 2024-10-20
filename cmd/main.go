package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jiazhen-lin/flight-booking/internal/app"
	"github.com/jiazhen-lin/flight-booking/internal/service"
	"github.com/sirupsen/logrus"
)

type config struct {
	httpAddr string
}

func initConfig() (*config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("empty env HTTP_ADDR")
	}

	return &config{
		httpAddr: addr,
	}, nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	// init
	config, err := initConfig()
	if err != nil {
		logrus.Fatalf("init config: %v", err)
	}
	logrus.Infof("config: %+v", config)

	flightService := service.NewFlightService()
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
