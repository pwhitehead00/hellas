package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	moduleRegistry "github.com/ironhalo/hellas/internal/moduleRegistry"
)

func download(rg *gin.RouterGroup, mr moduleRegistry.Registry) {
	download := rg.Group("/modules")

	download.GET("/:namespace/:name/:provider/:version/download", func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		provider := c.Param("provider")
		version := c.Param("version")

		url := mr.Download(namespace, name, provider, version)

		c.Header("X-Terraform-Get", url)
		c.Status(http.StatusNoContent)
	})
}

func version(rg *gin.RouterGroup, mr moduleRegistry.Registry) {
	versions := rg.Group("/modules")

	versions.GET("/:namespace/:name/:provider/versions", func(c *gin.Context) {
		namespace := c.Param("namespace")
		provider := c.Param("provider")
		name := c.Param("name")

		v, err := mr.ListVersions(namespace, name, provider)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		} else {
			o := mr.Versions(namespace, name, provider, v)
			c.JSON(http.StatusOK, o)
		}
	})
}

func ModuleRegistryGroup(rg *gin.Engine, mr moduleRegistry.Registry) {
	v1 := rg.Group("/v1")
	version(v1, mr)
	download(v1, mr)
}
