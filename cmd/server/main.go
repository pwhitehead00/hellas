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
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// todo - proper config setup
	config := mr.Config{
		Registry: map[string]any{
			"github": mr.GithubConfig{
				Protocol: "https",
			},
		},
	}
	// var config mr.Config

	// data, err := os.ReadFile("configFile")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := yaml.Unmarshal(data, &config); err != nil {
	// 	log.Fatal(err)
	// }

	mux, err := mr.NewModuleRegistry(config)
	if err != nil {
		log.Fatal(err)
	}

	// todo - configure TLS with config
	srv, err := server.NewServer(mux, true, "/tls/tls.crt", "/tls/tls.key")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server startup failed: %s\n", err)
		}
	}()

	log.Println("serving ...")

	<-ctx.Done()
	log.Println("shutdown requested:", ctx.Err())
	stop()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown forced: %s", err)
	}

	log.Println("shutdown complete")
}
