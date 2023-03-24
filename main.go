package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	links "github.com/shortenUrl/links"
	"github.com/shortenUrl/shortUrl"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	host := os.Getenv("DB_HOST") // Docker Compose service name
	port := os.Getenv("DB_PORT") // Default MySQL port
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	createDatabase := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	db.Exec(createDatabase)
	usedb := fmt.Sprintf("USE %s", dbName)
	db.Exec(usedb)

	migrationError := db.AutoMigrate(links.Links{})
	if migrationError != nil {
		log.Fatal(migrationError)
	}

	r.POST("/shortUrl", func(c *gin.Context) {
		shortUrl.ShortenUrl(db, c)
	})
	r.GET("/:shortURL", func(c *gin.Context) {
		shortUrl.RedirectToOriginalUrl(db, c)
	})
	r.Run(":5000")
}
