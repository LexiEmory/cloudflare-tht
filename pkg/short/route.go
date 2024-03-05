package short

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(conn *gorm.DB, router *gin.RouterGroup) {
	rg := router.Group("/s", RequestLogger(conn))

	rg.GET("/:id", Redirection(conn))
	rg.POST("/", CreateShort(conn))
	rg.DELETE("/:id", DeleteShort(conn))
	rg.GET("/:id/details", ShortDetails(conn))
}
