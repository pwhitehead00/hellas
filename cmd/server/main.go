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

func setupRouter(moduleType string, mr models.ModuleRegistry) *gin.Engine {
	registry := moduleregistry.NewModuleRegistry(moduleType, mr)
	r := gin.New()

	c := gin.LoggerConfig{
		SkipPaths: []string{"/healthcheck", "/.well-known/terraform.json"},
		Formatter: gin.LogFormatter(func(param gin.LogFormatterParams) string {
			return logging.Logger(param)
		}),
	}
	r.Use(gin.LoggerWithConfig(c), gin.Recovery())

	v1.ModuleRegistryGroup(r, registry)
	api.HealthCheck(r)
	api.WellKnown(r)

	return r
}

func main() {
	moduleType, ok := os.LookupEnv("MODULE_REGISTRY_TYPE")
	if !ok {
		log.Fatal("MODULE_REGISTRY_TYPE not set")
	}

	mrConfig, err := newMRConfig("/config/config.json")
	if err != nil {
		panic(err)
	}

	r := setupRouter(moduleType, *mrConfig)
	r.RunTLS(":8443", "/tls/tls.crt", "/tls/tls.key")
}
