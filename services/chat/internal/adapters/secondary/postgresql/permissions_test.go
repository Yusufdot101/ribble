package postgresql

import "github.com/Yusufdot101/ripple/services/chat/internal/application/core/domain"

func (rts *RepositoryTestSuite) TestNewPermission() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	permission := domain.NewPermission(domain.SendMessage)

	err = adapater.NewPermission(permission)
	rts.Require().Nil(err)
}

func (rts *RepositoryTestSuite) TestGetUserPermissions() {
	adapater, err := NewAdapter(rts.dataSourceURL)
	rts.Require().Nil(err)

	// create chat
	chat := domain.NewChat("", false)
	err = adapater.InsertChat(chat)
	rts.Require().Nil(err)

	// create role
	role := domain.NewRole(domain.Admin)
	err = adapater.NewRole(role)
	rts.Require().Nil(err)

	// create permission
	permission := domain.NewPermission(domain.SendMessage)
	err = adapater.NewPermission(permission)
	rts.Require().Nil(err)

	// create chat participant
	chatParticipant := domain.NewChatParticipant(1, chat.ID)
	err = adapater.InsertChatParticipant(chatParticipant)
	rts.Nil(err)

	// create chat role
	chatRole := domain.NewChatRole(chat.ID)
	err = adapater.NewChatRole(chatRole, role.Name)
	rts.Require().Nil(err)

	// grant permission to role
	err = adapater.GrantChatRolePermission(chatRole.ID, permission.Name)
	rts.Require().Nil(err)

	// grant chat role to user
	err = adapater.GrantUserChatRole(chatParticipant.UserID, chat.ID, role.Name)
	rts.Require().Nil(err)

	// get the permissions
	gotPermissions, err := adapater.GetUserPermissions(chatParticipant.UserID, chat.ID)
	rts.Require().Nil(err)

	rts.Require().Equal(1, len(gotPermissions))
	rts.Require().Equal(permission.ID, gotPermissions[0].ID)
	rts.Require().Equal(permission.Name, gotPermissions[0].Name)
}
