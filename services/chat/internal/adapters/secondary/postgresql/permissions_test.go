package postgresql

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

func (rts *RepositoryTestSuite) TestNewRole() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	role := domain.NewRole(domain.Admin)

	err = adapater.NewRole(role)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestNewPermission() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	permission := domain.NewPermission(domain.ReadMessage)

	err = adapater.NewPermission(permission)
	rts.Require().Nil(err)
}
