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

	var userID uint = 1
	message := domain.NewMessage(chat.ID, userID, "test message")
	err = adapater.InsertMessage(message)
	rts.Require().Nil(err)

	chatID, err := adapater.DeleteMessage(userID, message.ID)
	rts.Require().Nil(err)
	rts.Require().Equal(chat.ID, chatID)

	messages, err := adapater.GetMessages(chat.ID)
	rts.Require().Nil(err)
	rts.Require().Len(messages, 1)
	rts.Require().Equal("", messages[0].Content)
	rts.Require().True(messages[0].Deleted)
	rts.Require().False(messages[0].DeletedAt.IsZero())
}

func (rts *RepositoryTestSuite) TestEditMessage() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	var userID uint = 1
	message := domain.NewMessage(chat.ID, userID, "test message")
	err = adapater.InsertMessage(message)
	rts.Require().Nil(err)

	message, err = adapater.EditMessage(userID, message.ID, "new content")
	rts.Require().Nil(err)
	rts.Require().Equal(chat.ID, message.ChatID)

	messages, err := adapater.GetMessages(chat.ID)
	rts.Require().Nil(err)

	rts.Require().Equal("new content", messages[0].Content)
}

func (rts *RepositoryTestSuite) TestEditMessageFail() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	var userID uint = 1
	message := domain.NewMessage(chat.ID, userID, "original content")
	err = adapater.InsertMessage(message)
	rts.Require().Nil(err)

	message, err = adapater.EditMessage(2, message.ID, "new content")
	rts.Require().Error(err)
	rts.Require().Nil(message)

	messages, err := adapater.GetMessages(chat.ID)
	rts.Require().Nil(err)

	rts.Require().Equal("original content", messages[0].Content)
}
