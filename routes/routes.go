package routes

import (
	"fmt"

	"github.com/Er-Sadiq/controllers"
	"github.com/Er-Sadiq/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) {
	fmt.Println("🔧 Registering user routes...")

	route := r.Group("/api")
	{
		// Public routes
		route.POST("/register", func(c *gin.Context) {
			controllers.Register(c, db)
		})
		route.POST("/login", func(c *gin.Context) {
			controllers.Login(c, db)
		})

		// Send query route - Making it protected
		protected := route.Group("/")
		protected.Use(middleware.JWTAuthMiddleware()) // Protect this route with JWT middleware
		protected.POST("/send", controllers.SendQuery)
		protected.POST("/savequery", func(c *gin.Context) { controllers.SaveQuery(c, db) })
		protected.GET("/getqueries", func(c *gin.Context) { controllers.GetSavedQueries(c, db) })
		protected.DELETE("/deletequery/:id", func(ctx *gin.Context) {
			controllers.DeleteQuery(ctx, db)
		})
	}
}
