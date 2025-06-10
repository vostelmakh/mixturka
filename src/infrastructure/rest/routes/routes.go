package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/adapter"
)

func ApplicationRouter(router *gin.Engine, db *gorm.DB) {
	v1 := router.Group("/v1")

	AuthRoutes(v1, adapter.AuthAdapter(db))
	UserRoutes(v1, adapter.UserAdapter(db))
	MedicineRoutes(v1, adapter.MedicineAdapter(db))
}
