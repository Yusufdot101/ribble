package postgresql

import "github.com/Yusufdot101/ribble/services/chat/internal/application/core/domain"

func (rts *RepositoryTestSuite) TestInsertChat() {
	adapter, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)
	chat := domain.NewChat()
	err = adapter.InsertChat(chat)
	rts.Nil(err)
}
