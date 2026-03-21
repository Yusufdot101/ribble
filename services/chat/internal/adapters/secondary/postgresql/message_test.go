package postgresql

import (
	"github.com/Yusufdot101/ribble/services/chat/internal/application/core/domain"
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
