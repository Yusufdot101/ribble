package domain

type Chat struct {
	ID uint
}

func NewChat() *Chat {
	return &Chat{}
}
