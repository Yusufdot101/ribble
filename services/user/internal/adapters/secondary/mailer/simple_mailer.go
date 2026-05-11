package mailer

import (
	"log"
)

type Mailer struct{}

func NewMailer() *Mailer {
	return &Mailer{}
}

func (m *Mailer) Send(receipiant, email string) error {
	log.Println("receipiant: ", receipiant)
	log.Println("email: ", email)
	return nil
}
