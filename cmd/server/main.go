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

const (
	discoveryPath string = "GET /.well-known/terraform.json"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	mux := http.NewServeMux()
	mux.HandleFunc(discoveryPath, server.WellKnown)

	registry, err := mr.NewModuleRegistry(config)
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /github/{group}/{project}/versions", registry.Versions)
	mux.HandleFunc("GET /github/{group}/{project}/{version}/download", registry.Download)

	srv, err := server.NewServer(mux, false, "", "")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		// if err := srv.ListenAndServeTLS("", ""); err != nil && errors.Is(err, http.ErrServerClosed) {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
