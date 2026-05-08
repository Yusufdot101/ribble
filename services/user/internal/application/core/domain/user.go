package domain

import (
	"time"
)

type UserInfo struct {
	Provider string
	Sub      string
	Email    string
	Name     string
}

type UserIdentity struct {
	ID        uint
	Provider  string
	Sub       string
	UserID    uint
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewIdentity(provider, sub string) *UserIdentity {
	return &UserIdentity{
		Provider: provider,
		Sub:      sub,
	}
}

func NewUser(name, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}
