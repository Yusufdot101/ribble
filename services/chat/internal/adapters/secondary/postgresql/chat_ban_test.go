package postgresql

import (
	"time"

	"github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertChatBan() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat("", false)
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	expiration := time.Now().Add(time.Hour)
	chatBan := domain.NewChatBan(chat.ID, 1, 1, "test", &expiration)
	err = adapater.InsertChatBan(chatBan)
	rts.Nil(err)
	rts.Require().Equal(chat.ID, chatBan.ChatID)
}

func (rts *RepositoryTestSuite) TestGetChatBans() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat("", false)
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	expiration := time.Now().Add(time.Hour)
	chatBan := domain.NewChatBan(chat.ID, 1, 1, "", &expiration)
	err = adapater.InsertChatBan(chatBan)
	rts.Nil(err)
	rts.Require().Equal(chat.ID, chatBan.ChatID)

	chatBans, err := adapater.GetChatBans(chat.ID)
	rts.Require().Nil(err)
	rts.Require().Len(chatBans, 1)
	rts.Require().Equal(chatBan.ID, chatBans[0].ID)
	rts.Require().Equal(chatBan.UserID, chatBans[0].UserID)
}
