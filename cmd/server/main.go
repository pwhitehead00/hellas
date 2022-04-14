package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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

// Create a new config file
func newConfig(c []byte) (*models.Config, error) {
	var config models.Config

	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}

	config.ModuleBackend = strings.ToLower(config.ModuleBackend)
	config.ModuleRegistry.Protocol = strings.ToLower(config.ModuleRegistry.Protocol)

	if config.ModuleRegistry.Protocol != "https" && config.ModuleRegistry.Protocol != "ssh" {
		return nil, errors.New(fmt.Sprintf("Invalid protocol: %s", config.ModuleRegistry.Protocol))
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
	file, err := reader("/config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	c, err := newConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	r := newRouter(c)
	r.RunTLS(":8443", "/tls/tls.crt", "/tls/tls.key")
}
