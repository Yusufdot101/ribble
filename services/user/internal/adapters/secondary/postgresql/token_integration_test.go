package postgresql

import (
	"time"

	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
)

func (rts *RepositoryTestSuite) TestInsertToken() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)

	// insert user becauset the tokens table has reference to users id
	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	token := domain.NewToken(domain.REFRESH, domain.UUID, 1, "refreshToken", time.Now())
	err = adapter.InsertToken(token)
	rts.Nil(err)
}

func (rts *RepositoryTestSuite) TestGetTokenByStringAndUse() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)

	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	token := domain.NewToken(domain.REFRESH, domain.UUID, 1, "refreshToken", time.Now().Add(time.Hour))
	err = adapter.InsertToken(token)
	rts.Require().Nil(err)

	gotToken, err := adapter.GetTokenByStringAndUse(token.TokenString, token.Use)
	rts.Require().Nil(err)
	rts.Assert().Equal(token.TokenString, gotToken.TokenString)
	rts.Assert().Equal(token.Use, gotToken.Use)
	rts.Assert().Equal(token.TokenType, gotToken.TokenType)
	rts.Assert().Equal(token.UserID, gotToken.UserID)
}

func (rts *RepositoryTestSuite) TestDeleteTokenByStringAndUse() {
	adapter, err := NewAdapter(rts.DataSourceURL)
	rts.Require().Nil(err)

	user := domain.NewUser("yusuf", "example@gmail.com")
	err = adapter.InsertUser(user)
	rts.Require().Nil(err)

	token := domain.NewToken(domain.REFRESH, domain.UUID, user.ID, "refreshToken", time.Now().Add(time.Hour))
	err = adapter.InsertToken(token)
	rts.Require().Nil(err)

	err = adapter.DeleteTokenByStringAndUse(token.TokenString, token.Use)
	rts.Require().Nil(err)

	_, err = adapter.GetTokenByStringAndUse(token.TokenString, token.Use)
	rts.Require().Equal(domain.ErrRecordNotFound, err)
}
