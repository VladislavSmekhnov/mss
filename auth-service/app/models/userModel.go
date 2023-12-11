package models

import (
	"gorm.io/gorm"
)

type UserType string

const (
	Admin      UserType = "admin"
	Editor     UserType = "editor"
	Subscriber UserType = "subscriber"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null;default:null"`
	Password string `gorm:"not null;default:null"`
	Type     UserType
}
