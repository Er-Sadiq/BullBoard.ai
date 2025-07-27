package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique"`
	Password     string
	IsVerified   bool
	SavedQueries string `gorm:"type:text"`
}

type Query struct {
	ID       string `json:"id"`
	Query    string `json:"query"`
	Response string `json:"response"`
}
