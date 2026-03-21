package postgresql

import (
	"github.com/Yusufdot101/ribble/services/chat/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertChatParticipant() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipant(chatParticipant)
	rts.Nil(err)
}
