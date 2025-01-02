package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwhitehead00/hellas/internal/logging"
	mr "github.com/pwhitehead00/hellas/internal/moduleRegistry"
	"github.com/pwhitehead00/hellas/internal/server"
	"gopkg.in/yaml.v3"
)

var (
	logLevel string
)

func main() {
	flag.StringVar(&logLevel, "log-level", "debug", "set the logging level")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logging.Level.Set(logging.SetLogLevel(logLevel))

	configBytes, err := os.ReadFile("/config/config.yaml")
	if err != nil {
		logging.Log.Error("failed to read config", "error", err)
		os.Exit(1)
	}

	var config mr.Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		logging.Log.Error("failed to unmarshal config", "error", err)
		os.Exit(1)
	}

	mux, err := mr.NewModuleRegistry(config)
	if err != nil {
		logging.Log.Error("failed to create module registry", "error", err)
		os.Exit(1)
	}

	srv, err := server.NewServer(mux, true, "/tls/tls.crt", "/tls/tls.key")
	if err != nil {
		logging.Log.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			logging.Log.Error("server startup failed", "error", err)
			os.Exit(1)
		}
	}()

	logging.Log.Info("serving")

	<-ctx.Done()
	log.Println("shutdown requested:", ctx.Err())

	if err := srv.Shutdown(ctx); err != nil {
		logging.Log.Error("server shutdown forced", "error", err)
		os.Exit(1)
	}

	logging.Log.Info("shutdown complete")
}
