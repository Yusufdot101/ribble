package postgresql

import (
	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertChatParticipants() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat("", false)
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipants([]*domain.ChatParticipant{chatParticipant})
	rts.Nil(err)
}

func (rts *RepositoryTestSuite) TestDeleteChatParticipant() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat("", false)
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipants([]*domain.ChatParticipant{chatParticipant})
	rts.Nil(err)

	err = adapater.DeleteChatParticipant(chatParticipant.ChatID, chatParticipant.UserID)
	rts.Require().Nil(err)

	_, err = adapater.GetChatUsers(chatParticipant.ChatID, chatParticipant.UserID)
	rts.Require().Error(err)
	rts.Require().Equal(domain.ErrRecordNotFound, err)
}

func (rts *RepositoryTestSuite) TestGetChatUsers() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat("", false)
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipants([]*domain.ChatParticipant{chatParticipant})
	rts.Require().Nil(err)

	chatParticipant2 := domain.NewChatParticipant(2, chat.ID)
	err = adapater.InsertChatParticipants([]*domain.ChatParticipant{chatParticipant2})
	rts.Require().Nil(err)

	gotParticipants, err := adapater.GetChatUsers(chat.ID, chatParticipant.UserID)
	rts.Require().Nil(err)

	rts.True(len(gotParticipants) == 2)

	rts.True(chatParticipant.UserID == gotParticipants[0].UserID)
	rts.True(chatParticipant2.UserID == gotParticipants[1].UserID)
}
