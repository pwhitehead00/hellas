package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ironhalo/hellas/models"
)

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/.well-known/terraform.json", func(c *gin.Context) {
		var wk models.WellKnown
		wk.Modules = "/v1/modules/"
		c.JSON(200, wk)
	})

	r.GET("/v1/modules/terraform-aws-modules/vpc/aws/versions", func(c *gin.Context) {
		var m models.ModuleVersion
		var v models.Versions

		m.Modules = append(m.Modules, v)
		m.Modules[0].Versions = append(m.Modules[0].Versions, models.Version{"1.0.0"})
		m.Modules[0].Versions = append(m.Modules[0].Versions, models.Version{"1.1.0"})
		m.Modules[0].Versions = append(m.Modules[0].Versions, models.Version{"2.0.0"})

		c.JSON(200, m)
	})

	r.GET("/v1/modules/terraform-aws-modules/vpc/aws/1.0.0/download", func(c *gin.Context) {

		c.Header("X-Terraform-Get", "https://api.github.com/repos/hashicorp/terraform-aws-consul/tarball/v0.0.1//*?archive=tar.gz")
		c.Status(204)

	})

	r.RunTLS(":8080", "/app/server.crt", "/app/server.key") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
