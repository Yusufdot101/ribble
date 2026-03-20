package postgresql

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatID uint
	UserID uint
}
