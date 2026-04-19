package postgresql

import (
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertMessage() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	message := domain.NewMessage(chat.ID, 1, "test message")
	err = adapater.InsertMessage(message)
	rts.Nil(err)
}

func (rts *RepositoryTestSuite) TestGetMessages() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	message := domain.NewMessage(chat.ID, 1, "test message")
	err = adapater.InsertMessage(message)
	rts.Require().Nil(err)

	message2 := domain.NewMessage(chat.ID, 2, "test message reply")
	err = adapater.InsertMessage(message2)
	rts.Require().Nil(err)

	messages, err := adapater.GetMessages(chat.ID)
	rts.Require().Nil(err)

	rts.Require().Equal(2, len(messages))
	rts.Require().Equal(message.ID, messages[0].ID)
	rts.Require().Equal(message2.ID, messages[1].ID)
}

func (rts *RepositoryTestSuite) TestDeleteMessage() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	message := domain.NewMessage(chat.ID, 1, "test message")
	err = adapater.InsertMessage(message)
	rts.Require().Nil(err)

	err = adapater.DeleteMessage(message.ID)
	rts.Require().Nil(err)

	messages, err := adapater.GetMessages(chat.ID)
	rts.Require().Nil(err)
	rts.Require().Equal(0, len(messages))
}
