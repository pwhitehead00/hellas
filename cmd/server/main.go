package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwhitehead00/hellas/internal/server"
)

const (
	discoveryPath string = "GET /.well-known/terraform.json"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc(discoveryPath, server.WellKnown)

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
