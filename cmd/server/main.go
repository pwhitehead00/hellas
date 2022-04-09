package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	api "github.com/ironhalo/hellas/internal/api"
	v1 "github.com/ironhalo/hellas/internal/api/v1"
	"github.com/ironhalo/hellas/internal/logging"
	"github.com/ironhalo/hellas/internal/models"
	moduleregistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

func newMRConfig(c string) (*models.ModuleRegistry, error) {
	var mr models.ModuleRegistry

	file, err := os.ReadFile(c)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &mr)
	if err != nil {
		return nil, err
	}

	return &mr, nil
}

func newConfig(c string) (*models.Config, error) {
	var config models.Config

	file, err := os.ReadFile(c)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func newRouter(c *models.Config) *gin.Engine {
	registry := moduleregistry.NewModuleRegistry(c.ModuleBackend, *c.ModuleRegistry)
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

	return r
}

func main() {
	c, err := newConfig("/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	r := newRouter(c)

	r.RunTLS(":8443", "/tls/tls.crt", "/tls/tls.key")
}
