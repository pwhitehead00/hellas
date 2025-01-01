package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mr "github.com/pwhitehead00/hellas/internal/moduleRegistry"
	"github.com/pwhitehead00/hellas/internal/server"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	configBytes, err := os.ReadFile("/config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config mr.Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		log.Fatal(err)
	}

	mux, err := mr.NewModuleRegistry(config)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(mux, true, "/tls/tls.crt", "/tls/tls.key")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server startup failed: %s\n", err)
		}
	}()

	log.Println("serving ...")

	<-ctx.Done()
	log.Println("shutdown requested:", ctx.Err())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown forced: %s", err)
	}

	log.Println("shutdown complete")
}
