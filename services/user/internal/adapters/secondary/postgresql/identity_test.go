package postgresql

import (
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertIdentity() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)

	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	identity := domain.NewIdentity("local", "1")
	identity.UserID = user.ID
	err = adapter.InsertIdentity(identity)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) FindIdentityByProviderAndSub() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)

	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	identity := domain.NewIdentity("local", user.Email)
	identity.UserID = user.ID
	err = adapter.InsertIdentity(identity)
	rts.Require().Nil(err)

	gotIdentity, err := adapter.FindIdentityByProviderAndSub(identity.Provider, identity.Sub)
	rts.Require().Nil(err)
	rts.Require().Equal(identity.Provider, gotIdentity.Provider)
	rts.Require().Equal(identity.Sub, gotIdentity.Sub)
	rts.Require().Equal(identity.UserID, gotIdentity.UserID)
}
