package domain

type Message struct {
	ID       uint
	ChatID   uint
	SenderID uint
	Content  string
}
