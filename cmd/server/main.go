package main

import (
	"context"
	"encoding/json"
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
			"github": mr.GithubConfig{},
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
	mux.HandleFunc("/github/{group}/{project}", func(w http.ResponseWriter, r *http.Request) {
		group := r.PathValue("group")
		project := r.PathValue("project")
		data, err := registry.ListVersions(group, project)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("json encoding failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

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
