package postgresql

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ChatID uint
	UserID uint
}

type Chat struct {
	gorm.Model
	Users []User `gorm:"constraint:OnDelete:CASCADE;"`
}
