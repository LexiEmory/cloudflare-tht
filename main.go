package main

// Basic main boilerplate

import (
	"cloudflare-tht/pkg/short"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// ignore env loading errors, assume env is already configured in docker. its here simply for dev configuration
	godotenv.Load()

	// this Pings the database trying to connect
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%v port=%v dbname=%v password=%v username=%v sslmode=disable", os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDBNAME"), os.Getenv("PGPASSWORD"), os.Getenv("PGUSER")))
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	short.Route(db, r.Group("/"))

	r.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("HOST_PORT")))
}
