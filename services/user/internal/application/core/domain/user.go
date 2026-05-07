package domain

import (
	"time"
)

type Entry struct {
	ID       uint
	Sub      string
	Provider string
	UserID   uint
	Email    string
}

type User struct {
	ID        uint      `json:"id"` // this is local id
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(name, email, provider, sub string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Provider: provider,
		Sub:      sub,
	}
}
