package postgresql

import (
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertUser() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	saveErr := adapter.InsertUser(&domain.User{})
	rts.Nil(saveErr)
}

func (rts *RepositoryTestSuite) TestFindUserByEmail() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Nil(err)

	gotUser, err := adapter.FindUserByEmail(user.Email)
	rts.Nil(err)
	rts.Equal("example@gmail.com", gotUser.Email)
}
