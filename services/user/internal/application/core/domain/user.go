package domain

import (
	"time"
)

type UserInfo struct {
	Provider      string
	Sub           string
	Email         string
	Name          string
	EmailVerified bool
	PasswordHash  *[]byte
}

type UserIdentity struct {
	ID            uint
	Provider      string
	Sub           string
	UserID        uint
	EmailVerified bool
	PasswordHash  *[]byte
	CreatedAt     time.Time
	Email         string
	Name          string
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
