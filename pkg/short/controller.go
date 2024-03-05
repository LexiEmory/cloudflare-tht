package short

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Redirection handles redirecting and logging user activity when an interaction with a short url occurs
func Redirection(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		reqShort := Short{ID: id}
		res := conn.First(&reqShort)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			}
			return
		}

		log := AccessLog{
			ShortID:    reqShort.ID,
			AccessTime: time.Now().UTC(),
		}
		res = conn.Create(&log)
		if res.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			return
		}

		c.Redirect(http.StatusFound, reqShort.URL)
	}
}

// CreateShort handles creating a new short from a given url
func CreateShort(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ShortRequest
		err := c.BindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, JsonError(err))
			return
		}

		shortRecord := Short{URL: req.URL}
		res := conn.Create(&shortRecord)
		if res.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			return
		}

		c.JSON(http.StatusCreated, shortRecord)
	}
}

// DeleteShort handles deleting short urls
func DeleteShort(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		reqShort := Short{ID: id}
		res := conn.Delete(&reqShort)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ShortDetails handles gathering metrics on a short link sending it back to the requester
func ShortDetails(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		reqShort := Short{ID: id}
		res := conn.First(&reqShort)
		if res.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			return
		}

		response, err := BuildMetrics(conn, reqShort)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			return
		}

		c.JSON(http.StatusOK, *response)
	}
}
