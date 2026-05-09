package postgresql

import (
	"context"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertUser() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	saveErr := adapter.InsertUser(&domain.User{})
	rts.Require().Nil(saveErr)
}

func (rts *RepositoryTestSuite) TestFindUserByEmail() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	gotUser, err := adapter.FindUserByEmail(user.Email)
	rts.Require().Nil(err)
	rts.Require().Equal(user.Email, gotUser.Email)
}

func (rts *RepositoryTestSuite) TestFindUsersByID() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	ctx := context.Background()
	gotUsers, err := adapter.FindUsersByID(ctx, []uint32{uint32(user.ID)})
	rts.Require().Nil(err)
	rts.Require().Len(gotUsers, 1)
	rts.Require().Equal(user.ID, gotUsers[0].ID)
	rts.Require().Equal(user.Email, gotUsers[0].Email)
}

func (rts *RepositoryTestSuite) TestSearchUsers() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	user2 := domain.NewUser("yusuf2", "example2@gmail.com")
	err = adapter.InsertUser(user2)
	rts.Require().Nil(err)

	query := "yusuf"
	ctx := context.Background()
	gotUsers, err := adapter.SearchUsers(ctx, query, []uint32{uint32(user.ID), uint32(user2.ID)})
	rts.Require().Nil(err)
	rts.Require().Len(gotUsers, 2)

	query = "yusuf2"
	gotUsers, err = adapter.SearchUsers(ctx, query, []uint32{uint32(user.ID), uint32(user2.ID)})
	rts.Require().Nil(err)
	rts.Require().Len(gotUsers, 1)
	rts.Require().Equal(user2.ID, gotUsers[0].ID)
	rts.Require().Equal(user2.Email, gotUsers[0].Email)
}

func (rts *RepositoryTestSuite) TestFindUsersByEmail() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	gotUsers, err := adapter.FindUsersByEmail(user.Email)
	rts.Require().Nil(err)
	rts.Require().Len(gotUsers, 1)
	rts.Require().Equal(user.ID, gotUsers[0].ID)
	rts.Require().Equal(user.Email, gotUsers[0].Email)
}

func (rts *RepositoryTestSuite) TestGetContacts() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	user2 := domain.NewUser("yusuf2", "example2@gmail.com")
	err = adapter.InsertUser(user2)
	rts.Require().Nil(err)

	user3 := domain.NewUser("yusuf3", "example3@gmail.com")
	err = adapter.InsertUser(user3)
	rts.Require().Nil(err)

	ctx := context.Background()
	gotUsers, err := adapter.GetContacts(ctx, "", []uint32{}, uint32(user.ID))
	rts.Require().Nil(err)
	rts.Require().Len(gotUsers, 2)
	rts.Require().NotContains(gotUsers, user)
}
