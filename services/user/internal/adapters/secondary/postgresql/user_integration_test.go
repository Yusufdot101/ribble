package postgresql

import (
	"github.com/Yusufdot101/ribble/services/user/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertUser() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	saveErr := adapter.InsertUser(&domain.User{})
	rts.Nil(saveErr)
}

func (rts *RepositoryTestSuite) TestFindUserByProviderAndSub() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com", "google", "1")
	adapter.InsertUser(user)

	gotUser, err := adapter.FindUserByProviderAndSub(user.Provider, user.Sub)
	rts.Nil(err)
	rts.Equal("1", gotUser.Sub)
	rts.Equal("google", gotUser.Provider)
}
