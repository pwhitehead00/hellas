package main

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ironhalo/hellas/api/v1"
	moduleregistry "github.com/ironhalo/hellas/internal/moduleRegistry"
	"github.com/ironhalo/hellas/models"
)

func main() {

	registry := moduleregistry.NewModuleRegistry("github")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1.ModuleRegistryGroup(r, registry)

	r.GET("/.well-known/terraform.json", func(c *gin.Context) {
		var wk models.WellKnown
		wk.Modules = "/v1/modules/"
		c.JSON(200, wk)
	})

	r.RunTLS(":8080", "/app/server.crt", "/app/server.key") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
