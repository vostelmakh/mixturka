package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/adapter"
)

func ApplicationRouter(router *gin.Engine, db *gorm.DB) {
	v1 := router.Group("/v1")

	v1.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "1.0.1",
			"service": "Mixturka",
		})
	})

	AuthRoutes(v1, adapter.AuthAdapter(db))
	UserRoutes(v1, adapter.UserAdapter(db))
	MedicineRoutes(v1, adapter.MedicineAdapter(db))
}
