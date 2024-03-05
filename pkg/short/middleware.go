package short

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func RequestLogger(conn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// pre request
		startTime := time.Now().UTC()

		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))

		audit := Audit{
			Action:  c.Request.Method,
			Path:    c.Request.URL.Path,
			Request: string(jsonData[:]),
			Invoker: c.ClientIP(),
		}

		c.Next()

		audit.Latency = time.Since(startTime).String()

		res := conn.Create(&audit)
		if res.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, JsonError(res.Error))
			return
		}
	}
}
