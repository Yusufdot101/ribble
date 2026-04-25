package postgresql

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

func (rts *RepositoryTestSuite) TestNewRole() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	role := domain.NewRole(domain.Admin)

	err = adapater.NewRole(role)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestNewChatRole() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create role
	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	chatRole := domain.NewChatRole(chat.ID)
	err = adapater.NewChatRole(chatRole, role.Name)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestNewChatRoleFail() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// invalid role id
	chatRole := domain.NewChatRole(chat.ID)
	err = adapater.NewChatRole(chatRole, domain.Admin)
	rts.Require().NotNil(err)
	rts.Require().Equal(domain.ErrInvalidChatRole, err)
}

func (rts *RepositoryTestSuite) TestGrantChatRolePermission() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create role
	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// create permission
	permission := domain.NewPermission(domain.SendMessage)
	err = adapater.NewPermission(permission)
	rts.Require().Nil(err)

	// create chat role
	chatRole := domain.NewChatRole(chat.ID)
	err = adapater.NewChatRole(chatRole, role.Name)
	rts.Require().Nil(err)

	// grant permission to chat role
	err = adapater.GrantChatRolePermission(chatRole.ID, permission.Name)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestGrantUserChatRole() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create permission
	permission := domain.NewPermission(domain.SendMessage)
	err = adapater.NewPermission(permission)
	rts.Require().Nil(err)

	// create role
	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat()
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// create chat role
	chatRole := domain.NewChatRole(chat.ID)
	err = adapater.NewChatRole(chatRole, role.Name)
	rts.Require().Nil(err)

	// create chat participant
	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipant(chatParticipant)
	rts.Nil(err)

	err = adapater.GrantUserChatRole(chatParticipant.UserID, chat.ID, domain.Admin)
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
	err = adapater.GrantUserChatRole(chatParticipant.UserID, chat.ID, domain.Admin)
	rts.Require().Equal(domain.ErrInvalidRole, err)
}
