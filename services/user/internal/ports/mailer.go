package ports

type Mailer interface {
	Send(receipiant, email string) error
}
