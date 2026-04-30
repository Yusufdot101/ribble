package postgresql

import (
	"context"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertChat() {
	adapter, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)
	chat := domain.NewChat("", false)
	err = adapter.InsertChat(chat)
	rts.Nil(err)
}

func (rts *RepositoryTestSuite) TestGetChatByUserIDs() {
	adapter, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat("", false)
	err = adapter.InsertChat(chat)
	rts.Require().Nil(err)

	// add users
	userID1 := 1
	userID2 := 2

	participant1 := domain.NewChatParticipant(uint(userID1), chat.ID)
	participant2 := domain.NewChatParticipant(uint(userID2), chat.ID)
	err = adapter.InsertChatParticipants([]*domain.ChatParticipant{participant1, participant2})
	rts.Require().Nil(err)

	// fetch the chat using those users' ids
	gotChat, err := adapter.GetChatByUserIDs([]uint{participant1.ID, participant2.ID}, false)
	rts.Require().Nil(err)
	rts.Require().Equal(chat.ID, gotChat.ID)
}

func (rts *RepositoryTestSuite) TestGetChatsByUserID() {
	err := rts.truncateDB(context.Background())
	rts.Require().NoError(err)
	adapter, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat("", false)
	err = adapter.InsertChat(chat)
	rts.Require().Nil(err)

	chat2 := domain.NewChat("", false)
	err = adapter.InsertChat(chat2)
	rts.Require().Nil(err)

	// add users
	userID := 1

	participant1 := domain.NewChatParticipant(uint(userID), chat.ID)
	participant2 := domain.NewChatParticipant(uint(userID), chat2.ID)
	err = adapter.InsertChatParticipants([]*domain.ChatParticipant{participant1, participant2})
	rts.Require().Nil(err)

	// fetch the chat using those users' ids
	chats, err := adapter.GetChatsByUserID(uint(userID), "")
	rts.Require().Nil(err)
	rts.Require().Equal(2, len(chats))
}
