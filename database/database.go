package database

import (
	"fmt"
	"log"

	"github.com/Er-Sadiq/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	// dsn := "root:admin@tcp(127.0.0.1:3306)/bullb_db?charset=utf8mb4&parseTime=True&loc=Local"
	// dsn := "root:IyrdoNPUBYELCmmNziPTRlJxRKteZETX@tcp(shuttle.proxy.rlwy.net:20526)/railway"
	dsn := "root:FozIrHdWoLvYotfUsEJlnSSIYZCdeVXu@tcp(nozomi.proxy.rlwy.net:15516)/railway?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 🔧 Create table automatically
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}

	fmt.Println("Connected to database")
	fmt.Println("✅ Connected to Railway MySQL!")

	DB = db
	return db
}
