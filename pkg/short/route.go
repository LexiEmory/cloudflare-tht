package short

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Route(conn *gorm.DB, router *gin.RouterGroup) {
	rg := router.Group("/s")

	rg.GET("/:id", func(c *gin.Context) {

		c.Status(http.StatusOK)
	})
}
