package grpc

import (
	"context"

	"github.com/Yusufdot101/ripple-proto/golang/user"
)

func (a *Adapter) VerifyUsers(ctx context.Context, userIDs []uint) (bool, error) {
	userIDs32 := []uint32{}
	for _, userID := range userIDs {
		userIDs32 = append(userIDs32, uint32(userID))
	}

	req := &user.VerifyUsersRequest{
		Ids: userIDs32,
	}
	res, err := a.userClient.VerifyUsers(ctx, req)
	if err != nil {
		return false, err
	}
	return res.AllValid, nil
}
