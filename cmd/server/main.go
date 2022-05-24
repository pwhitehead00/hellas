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
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/ironhalo/hellas/internal/api"
	v1 "github.com/ironhalo/hellas/internal/api/v1"
	"github.com/ironhalo/hellas/internal/logging"
	moduleregistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

func reader(f string) ([]byte, error) {
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Create a new gin Engine
func newRouter(r moduleregistry.Registry) (*gin.Engine, error) {
	g := gin.New()
	l := gin.LoggerConfig{
		SkipPaths: []string{"/healthcheck", "/.well-known/terraform.json"},
		Formatter: gin.LogFormatter(func(param gin.LogFormatterParams) string {
			return logging.Logger(param)
		}),
	}
	g.Use(gin.LoggerWithConfig(l), gin.Recovery())

	v1.ModuleRegistryGroup(g, r)
	api.HealthCheck(g)
	api.WellKnown(g)

	return g, nil
}

func main() {
	mrBackend := flag.String("module-registry-backend", "", "Terraform module registry backend")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	mrConfig, err := reader("/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	registry, err := moduleregistry.NewModuleRegistry(mrBackend, mrConfig)
	if err != nil {
		log.Fatal(err)
	}

	router, err := newRouter(registry)
	if err != nil {
		log.Fatal(err)
	}

	router.SetTrustedProxies(nil)

	srv := &http.Server{
		Addr:    ":8443",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServeTLS("/tls/tls.crt", "/tls/tls.key"); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
