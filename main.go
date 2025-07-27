package main

import (
	"time"

	"github.com/Er-Sadiq/database"
	"github.com/Er-Sadiq/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.Connect()

	r := gin.Default()
	r.SetTrustedProxies(nil) // optional
	// r.Use(cors.Default())

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"}, // Vite default port
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // ✅ Allow requests from any domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Needed for cookies/auth
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// ✅ THIS IS CRITICAL — CALL THIS FUNCTION
	routes.UserRoutes(r, db)

	r.Run(":8080")
}
