package main

import (
	"github.com/gin-gonic/gin"
	root "github.com/ironhalo/hellas/api"
	v1 "github.com/ironhalo/hellas/api/v1"
	moduleregistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

func setupRouter() *gin.Engine {
	registry := moduleregistry.NewModuleRegistry("github")
	r := gin.Default()

	v1.ModuleRegistryGroup(r, registry)
	root.HealthCheck(r)
	root.WellKnown(r)

	return r
}

func main() {
	r := setupRouter()
	r.RunTLS(":8443", "/app/server.crt", "/app/server.key")
}
