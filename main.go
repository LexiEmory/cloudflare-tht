package main

// Basic main boilerplate

import (
	"cloudflare-tht/pkg/short"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// ignore env loading errors, assume env is already configured in docker. its here simply for dev configuration
	godotenv.Load()

	// this Pings the database trying to connect
	dsn := fmt.Sprintf("host=%v port=%v dbname=%v password=%v user=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PASS"), os.Getenv("DB_USER"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
		return
	}

	r := gin.Default()

	short.Route(db, r.Group("/"))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("HOST_PORT")))
}
