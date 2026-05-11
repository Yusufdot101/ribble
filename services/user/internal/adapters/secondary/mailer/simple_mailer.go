package mailer

import (
	"log"
)

type Mailer struct{}

func NewMailer() *Mailer {
	return &Mailer{}
}

func (m *Mailer) Send(recipient, email string) error {
	log.Println("recipient: ", recipient)
	log.Println("email: ", email)
	return nil
}
