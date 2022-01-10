package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	root "github.com/ironhalo/hellas/api"
	v1 "github.com/ironhalo/hellas/api/v1"
	moduleregistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

func setupRouter(moduleType string) *gin.Engine {
	registry := moduleregistry.NewModuleRegistry(moduleType)
	r := gin.Default()

	v1.ModuleRegistryGroup(r, registry)
	root.HealthCheck(r)
	root.WellKnown(r)

	return r
}

func main() {
	moduleType, ok := os.LookupEnv("MODULE_REGISTRY_TYPE")
	if !ok {
		log.Fatal("MODULE_REGISTRY_TYPE not set")
	}

	r := setupRouter(moduleType)
	r.RunTLS(":8443", "/app/server.crt", "/app/server.key")
}
