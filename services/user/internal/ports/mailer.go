package ports

type Mailer interface {
	Send(recipient, email string) error
}
