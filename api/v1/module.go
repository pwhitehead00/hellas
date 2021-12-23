package v1

import (
	"github.com/gin-gonic/gin"
	moduleRegistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

// https://registry.terraform.io/v1/modules/terraform-aws-modules/vpc/aws/3.11.0/download
func download(rg *gin.RouterGroup, mr moduleRegistry.Registry) {
	download := rg.Group("/modules")

	download.GET("/:namespace/:name/:provider/:version/download", func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		provider := c.Param("provider")
		version := c.Param("version")

		url := mr.Download(namespace, name, provider, version)
		// url := fmt.Sprintf("git::https://github.com/%s/terraform-%s-%s?ref=v%s", namespace, provider, name, version)

		c.Header("X-Terraform-Get", url)
		c.Status(204)
	})
}

// https://github.com/terraform-aws-modules/terraform-aws-vpc
// https://registry.terraform.io/v1/modules/terraform-aws-modules/vpc/aws/versions
func version(rg *gin.RouterGroup, mr moduleRegistry.Registry) {
	versions := rg.Group("/modules")

	versions.GET("/:namespace/:name/:provider/versions", func(c *gin.Context) {
		namespace := c.Param("namespace")
		provider := c.Param("provider")
		name := c.Param("name")

		o := mr.Versions(namespace, name, provider)
		// o := gh.MegaGitHub(namespace, provider, name)
		c.JSON(200, o)
	})
}

func ModuleRegistryGroup(rg *gin.Engine, mr moduleRegistry.Registry) {
	v1 := rg.Group("/v1")
	version(v1, mr)
	download(v1, mr)
}
