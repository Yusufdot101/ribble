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

func (rts *RepositoryTestSuite) TestGrantRolePermission() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	permission := domain.NewPermission(domain.ReadMessage)
	err = adapater.NewPermission(permission)
	rts.Require().Nil(err)

	err = adapater.GrantRolePermission(role.ID, permission.Name)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestGrantRolePermissionFail() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	// permission not in the database, should error
	err = adapater.GrantRolePermission(role.ID, domain.ReadMessage)
	rts.Require().Equal(domain.ErrInvalidPermission, err)
}

func (rts *RepositoryTestSuite) TestGrantUserRole() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// create chat participant
	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipant(chatParticipant)
	rts.Nil(err)

	// create role
	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	err = adapater.GrantUserRole(chatParticipant.ID, role.Name)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestGrantUserRoleFail() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// create chat participant
	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipant(chatParticipant)
	rts.Nil(err)

	// role not in the database, should error
	err = adapater.GrantUserRole(chatParticipant.ID, domain.Admin)
	rts.Require().Equal(domain.ErrInvalidRole, err)
}

func (rts *RepositoryTestSuite) TestGetUserRole() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// create chat participant
	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipant(chatParticipant)
	rts.Nil(err)

	// create role
	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	// grant role to user
	err = adapater.GrantUserRole(chatParticipant.ID, role.Name)
	rts.Require().Nil(err)

	// get the role
	gotRole, err := adapater.GetUserRole(chatParticipant.UserID)
	rts.Require().Nil(err)

	rts.Require().Equal(role.ID, gotRole.ID)
	rts.Require().Equal(role.Name, gotRole.Name)
}
