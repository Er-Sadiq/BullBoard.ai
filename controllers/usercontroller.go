package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Er-Sadiq/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Protected(c *gin.Context) {
	email, _ := c.Get("userEmail")
	c.JSON(http.StatusOK, gin.H{
		"message": "Protected access granted",
		"user":    email,
	})
}

func SaveQuery(c *gin.Context, db *gorm.DB) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email := emailVal.(string)

	var newQuery models.Query
	if err := c.ShouldBindJSON(&newQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var user models.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Parse existing SavedQueries (stored as JSON string)
	var existingQueries []models.Query
	if user.SavedQueries != "" {
		if err := json.Unmarshal([]byte(user.SavedQueries), &existingQueries); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse existing saved queries"})
			return
		}
	}

	// Append new query
	existingQueries = append(existingQueries, newQuery)

	// Marshal back to JSON string
	updatedJSON, err := json.Marshal(existingQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode updated saved queries"})
		return
	}

	// Save to DB
	user.SavedQueries = string(updatedJSON)
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated queries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query saved successfully"})
}

func GetSavedQueries(c *gin.Context, db *gorm.DB) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email := emailVal.(string)

	var user models.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var savedQueries []models.Query
	if user.SavedQueries != "" {
		if err := json.Unmarshal([]byte(user.SavedQueries), &savedQueries); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode saved queries"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"queries": savedQueries, // ✅ make sure this is an array
	})

}

func DeleteQuery(c *gin.Context, db *gorm.DB) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email := emailVal.(string)

	// Parse JSON body to get query ID
	var input models.Query
	if err := c.ShouldBindJSON(&input); err != nil || input.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Query ID required."})
		return
	}

	// Fetch user
	var user models.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Decode saved queries
	var savedQueries []models.Query
	if user.SavedQueries != "" {
		if err := json.Unmarshal([]byte(user.SavedQueries), &savedQueries); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse saved queries"})
			return
		}
	}

	// Filter out query with matching ID
	found := false
	filtered := make([]models.Query, 0)
	for _, q := range savedQueries {
		if q.ID == input.ID {
			found = true
			continue
		}
		filtered = append(filtered, q)
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Query not found"})
		return
	}

	// Encode and save updated list
	updatedJSON, err := json.Marshal(filtered)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode updated queries"})
		return
	}

	user.SavedQueries = string(updatedJSON)
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated queries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query deleted successfully"})
}
