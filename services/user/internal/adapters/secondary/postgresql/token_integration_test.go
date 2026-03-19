package postgresql

import (
	"time"

	"github.com/Yusufdot101/ribble/services/user/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertToken() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)

	token := domain.NewToken(domain.REFRESH, domain.UUID, 1, "refreshToken", time.Now())
	err = adapter.InsertToken(token)
	rts.Nil(err)
}
