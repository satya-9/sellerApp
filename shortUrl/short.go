package shortUrl

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	links "github.com/shortenUrl/links"
	"gorm.io/gorm"
)

func ShortenUrl(db *gorm.DB, c *gin.Context) {
	var requestBody struct {
		URL string `json:"url" binding:"required"`
	}
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Generate a unique short URL
	var shortURL string
	for {
		// Generate a random 6-byte token
		tokenBytes := make([]byte, 6)
		_, err := rand.Read(tokenBytes)
		if err != nil {
			log.Fatalf("Failed to generate random token: %v", err)
		}
		token := base64.RawURLEncoding.EncodeToString(tokenBytes)

		// Check if the token already exists in the database
		var count int64
		db.Model(links.Links{}).Where("short_url = ?", token).Count(&count)
		if count == 0 {
			shortURL = token
			break
		}
	}

	// Create a new URL object
	newURL := links.Links{
		OriginalUrl: requestBody.URL,
		ShortUrl:    shortURL,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	err = db.Create(&newURL).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	// Return the short URL to the client
	c.JSON(http.StatusOK, gin.H{"short_url": "https://my-short-link/" + shortURL})

}

func RedirectToOriginalUrl(db *gorm.DB, c *gin.Context) {
	type Result struct {
		ID          int
		ExpiresAt   string
		OriginalURL string
	}
	var result Result
	err := db.Model(&links.Links{}).Select("id,expires_at,original_url").Where("short_url = ?", c.Param("shortURL")).First(&result).Error
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	layout := "2006-01-02 15:04:05"
	expiresAtformat, err := time.Parse(layout, result.ExpiresAt)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}

	// Check if the short URL has expired
	if time.Now().After(expiresAtformat) {
		// Delete the expired URL object from the database
		db.Delete(&result.OriginalURL)

		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	c.Redirect(http.StatusMovedPermanently, result.OriginalURL)
}
