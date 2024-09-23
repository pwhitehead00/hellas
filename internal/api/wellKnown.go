package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pwhitehead00/hellas/internal/models"
)

func WellKnown(r *gin.Engine) {
	r.GET("/.well-known/terraform.json", func(c *gin.Context) {
		var wk models.WellKnown
		wk.Modules = "/v1/modules/"
		c.JSON(http.StatusOK, wk)
	})
}
