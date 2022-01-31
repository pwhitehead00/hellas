package root

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ironhalo/hellas/models"
)

func HealthCheck(r *gin.Engine) {
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
}

func WellKnown(r *gin.Engine) {
	r.GET("/.well-known/terraform.json", func(c *gin.Context) {
		var wk models.WellKnown
		wk.Modules = "/v1/modules/"
		c.JSON(http.StatusOK, wk)
	})
}
