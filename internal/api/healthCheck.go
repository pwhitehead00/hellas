package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(r *gin.Engine) {
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
}
