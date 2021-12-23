package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	gh "github.com/ironhalo/hellas/internal/github"
)

// https://registry.terraform.io/v1/modules/terraform-aws-modules/vpc/aws/3.11.0/download
func addDownload(rg *gin.RouterGroup) {
	download := rg.Group("/modules")

	download.GET("/:namespace/:name/:provider/:version/download", func(c *gin.Context) {
		namespace := c.Param("namespace")
		provider := c.Param("provider")
		name := c.Param("name")
		version := c.Param("version")
		url := fmt.Sprintf("git::https://github.com/%s/terraform-%s-%s?ref=v%s", namespace, provider, name, version)

		c.Header("X-Terraform-Get", url)
		c.Status(204)
	})
}

// https://github.com/terraform-aws-modules/terraform-aws-vpc
// https://registry.terraform.io/v1/modules/terraform-aws-modules/vpc/aws/versions
func addVersion(rg *gin.RouterGroup) {
	versions := rg.Group("/modules")

	versions.GET("/:namespace/:name/:provider/versions", func(c *gin.Context) {
		namespace := c.Param("namespace")
		provider := c.Param("provider")
		name := c.Param("name")

		o := gh.MegaGitHub(namespace, provider, name)
		c.JSON(200, o)
	})
}

func ModuleRegistry(rg *gin.Engine) {
	v1 := rg.Group("/v1")
	addVersion(v1)
	addDownload(v1)
}
