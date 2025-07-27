package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Er-Sadiq/models"
	"github.com/Er-Sadiq/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegVar struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const TavilyAPIKey = "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"

func Register(c *gin.Context, db *gorm.DB) {
	var regvar RegVar

	if err := c.ShouldBindJSON(&regvar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := utils.HashPassword(regvar.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     regvar.Name,
		Email:    regvar.Email,
		Password: hashedPassword,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or invalid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context, db *gorm.DB) {
	var auth Auth

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", auth.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPasswordHash(auth.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 🍪 Set HttpOnly cookie
	// c.SetCookie("jwt", token, 3600, "/", "localhost", false, true) // `true` for HttpOnly
	// c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func SendQuery(c *gin.Context) {
	var requestBody map[string]string

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	query, ok := requestBody["query"]
	if !ok || query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'query' field is required"})
		return
	}

	response, err := callTavilyAPI(query)
	if err != nil {
		log.Println("Error calling Tavily API:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from Tavily"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func callTavilyAPI(query string) (map[string]interface{}, error) {
	url := "https://api.tavily.com/search"

	payload := map[string]string{"query": query}
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+TavilyAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
