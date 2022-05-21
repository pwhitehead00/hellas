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
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/ironhalo/hellas/internal/api"
	v1 "github.com/ironhalo/hellas/internal/api/v1"
	"github.com/ironhalo/hellas/internal/logging"
	"github.com/ironhalo/hellas/internal/models"
	moduleregistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

func reader(f string) ([]byte, error) {
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Create and validate a new instance of models.Config
func newConfig(c []byte) (*models.Config, error) {
	var config models.Config

	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Create a new gin Engine
func newRouter(c *models.Config) (*gin.Engine, error) {
	registry, err := moduleregistry.NewModuleRegistry(c.ModuleBackend, *c.ModuleRegistry)
	if err != nil {
		return nil, err
	}

	r := gin.New()
	l := gin.LoggerConfig{
		SkipPaths: []string{"/healthcheck", "/.well-known/terraform.json"},
		Formatter: gin.LogFormatter(func(param gin.LogFormatterParams) string {
			return logging.Logger(param)
		}),
	}
	r.Use(gin.LoggerWithConfig(l), gin.Recovery())

	v1.ModuleRegistryGroup(r, registry)
	api.HealthCheck(r)
	api.WellKnown(r)

	return r, nil
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	file, err := reader("/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	config, err := newConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	router, err := newRouter(config)
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
