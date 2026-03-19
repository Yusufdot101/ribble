package domain

import (
	"time"
)

type User struct {
	ID        uint      `json:"id"`  // this is local id
	Sub       string    `json:"sub"` // this is id from providers(in case user changes email, the sub will be always be the same)
	Provider  string    `json:"provider"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
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
