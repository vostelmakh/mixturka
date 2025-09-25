package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplicationRouter(router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "1.0.1",
			"service": "Mixturka",
		})
	})
}
